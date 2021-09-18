package main

import (
	"os"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"log"

	"golang.org/x/crypto/bcrypt"
	"cloud.google.com/go/datastore"
)

// Player as it is stored in db.
type Player struct {
	Name  string
	Pwd   string
	Salt  string
	Score int64
	Sig   string
}

// PlayerScore returned to clients.
type PlayerScore struct {
	Name  string
	Score int64
}

const playerEntityKind string = "player"

var errPasswordInvalid = errors.New("Invalid password")
var errSigInvalid = errors.New("Invalid sig")

func main() {
	http.HandleFunc("/", handle)
	
	port := os.Getenv("PORT")
	if port == "" {
        port = "8080"
        log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	dsClient, err := datastore.NewClient(r.Context(), "homemaker-225521")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer dsClient.Close()

	if r.Method == "POST" {
		storeHighscore(dsClient, w, r)
	} else if r.Method == "GET" {
		getHighscores(dsClient, w, r)
	}
}

func getHighscores(dsClient *datastore.Client, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
 
	q := datastore.NewQuery(playerEntityKind).Order("-Score")
	var scores []PlayerScore
	if _, err := dsClient.GetAll(ctx, q, &scores); err != nil {
		res1B, _ := json.Marshal(scores)
		fmt.Fprintf(w, string(res1B))
	} else {
		log.Fatalf("Error: %v", err)
		fmt.Fprintf(w, "[]")
	}
}

func storeHighscore(dsClient *datastore.Client, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	player, err := parsePlayerRequestBody(w, r)
	if err != nil {
		http.Error(w, errors.New("Error parsing request body").Error(), http.StatusBadRequest)
		return
	}

	password := player.Pwd
	if password == "" {
		http.Error(w, errors.New("Empty password").Error(), http.StatusBadRequest)
		return
	}

	key := datastore.NameKey(playerEntityKind, player.Name, nil)

	storedPlayer := new(Player)
	_, err = dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		datastoreGetErr := tx.Get(key, storedPlayer)
		if datastoreGetErr != nil && datastoreGetErr != datastore.ErrNoSuchEntity {
			return datastoreGetErr
		}

		if checkPasswordCorrect(ctx, password, storedPlayer) {
			player.Salt = randomString(10)
			hash, err := hashPassword(ctx, password+player.Salt)
			if err == nil {
				player.Pwd = hash
				_, err = tx.Put(key, &player)
			}
		} else {
			err = errPasswordInvalid
		}

		if !checkSigCorrect(player) {
			err = errSigInvalid
		}

		return err
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
		if err == errPasswordInvalid {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func checkPasswordCorrect(c context.Context, password string, storedPlayer *Player) bool {
	if storedPlayer.Salt == "" {
		return true // No stored player
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedPlayer.Pwd), []byte(password+storedPlayer.Salt))
	return err == nil
}

func parsePlayerRequestBody(w http.ResponseWriter, r *http.Request) (Player, error) {
	player := Player{}
	err := json.NewDecoder(r.Body).Decode(&player)
	return player, err
}

func hashPassword(c context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	hash := string(bytes)
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	}
	return hash, err
}

func randomString(len int) string {
	buff := make([]byte, len)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	return str[:len]
}

func checkSigCorrect(p Player) bool {
	h := sha256.New()
	h.Write([]byte("hr" + p.Name + strconv.FormatInt(p.Score, 10) + "salt"))
	hexString := fmt.Sprintf("%x", h.Sum(nil))
	return strings.EqualFold(hexString, p.Sig)
}
