# GoScores
Leaderboard server on Golang, uses Cloud Datastore for storage.


#### Player JSON
{ "Name":"Player name","Score":101800 }

#### GET 
Returns JSON array of players
[{ "Name":"Player name","Score":101800 }]

#### POST
{ "Name":"Player name","Pwd":"Password","Score":11320, Sig: "signature" }

Note: Sig field is used for validation. 

#### Signature (JSON Sig field)
To prevent manual posting of scores this leaderboard uses signature. Signature is generated on client and
checked on the server, see function checkSigCorrect() in scores.go source file.
 
#### Leaderboard sample published app
https://play.google.com/store/apps/details?id=com.kurovsky.christmashomes

#### Android sample for posting highscores
``` 
	data class UserScore(val Name: String = "", val Score: Long? = 0L, val Sig: String = "")

        coroutineScope.launch {
            val userScore = UserScore("Name", 100, "Signature")
            val responseCode = URL("url")
                    .openConnection()
                    .let {
                        it as HttpURLConnection
                    }.run {
                        setRequestProperty("Content-Type", "application/json; charset=utf-8")
                        requestMethod = "POST"
                        doOutput = true
                        val outputWriter = OutputStreamWriter(outputStream)
                        outputWriter.write(Gson().toJson(userScore))
                        outputWriter.flush()
                        outputWriter.close()
                        responseCode
                    }
        }
```

Make sure to customize the sigString, and modify checkSigCorrect() function from scores.go accordingly. 


#### Android sample for getting highscores
```
	coroutineScope.launch {
            val usersScoresJson = URL("url").readText()
            val usersScores = Gson().fromJson<List<UserScore>>(usersScoresJson, object : TypeToken<List<UserScore>>() {}.type)
        }
```
   
