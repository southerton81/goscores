# GoScores
Leaderboard server on Golang, tailored to be hosted on App Engine. Uses Cloud Datastore for storage.


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
    class UserScore {
        @SerializedName("Name")
        public String Name = "";

        @SerializedName("Score")
        public Long Score = 0L;

        @SerializedName("Sig")
        public String Sig = "";
    }
    
    class PostScoreTask extends AsyncTask  {
        @Override
        protected Object doInBackground(Object[] strings) {
            BufferedReader in = null;
            Integer responseCode = null;

            try {
                URL obj = new URL(URL);
                HttpsURLConnection con = (HttpsURLConnection) obj.openConnection();
                con.setRequestMethod("POST");
                con.setRequestProperty("User-Agent", "app");
                con.setRequestProperty("Accept-Language", "en-US,en;q=0.5");
                String name = (String) strings[0];
                String pwd = (String) strings[1];
                String score = (String) strings[2];
                String urlParameters = "{\"Name\":\"" + name +
                        "\",\"Pwd\":\"" + pwd +
                        "\",\"Score\":" + score +
                        ",\"Sig\":\"" + generateSig(name, score) + "\"}";
                con.setDoOutput(true);
                DataOutputStream wr = new DataOutputStream(con.getOutputStream());
                wr.writeBytes(urlParameters);
                wr.flush();
                wr.close();
                responseCode = con.getResponseCode();
            } catch (Exception e) {
                e.printStackTrace();
            } finally {
                if (in != null) {
                    try {
                        in.close();
                    } catch (IOException e) {
                        e.printStackTrace();
                    }
                }

            }

            return responseCode;
        }

        @Override
        protected void onPostExecute(Object o) {
            super.onPostExecute(o);
            // ...
        }
    }

   public static String generateSig(String name, String score) {
        try {
            String sigString = "hr" + name + score + "salt";
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(sigString.getBytes("UTF-8"));
            return bytesToHex(hash);
        } catch (Exception e) {
            e.printStackTrace();
            return "congrats";
        }
    }

    private static String bytesToHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder();
        for (byte aByte : bytes) {
            String st = String.format("%02X", aByte);
            sb.append(st);
        }
        return sb.toString();
    }
```

Make sure to customize the sigString, and modify checkSigCorrect() function from scores.go accordingly. 


#### Android sample for getting highscores
```
 class GetScoresTask extends AsyncTask {
        @Override
        protected Object doInBackground(Object[] strings) {
            BufferedReader in = null;

            Integer responseCode = null;
            UserScore[] data = null;

            try {
                URL obj = new URL(URL);
                HttpsURLConnection con = (HttpsURLConnection) obj.openConnection();
                con.setRequestMethod("GET");
                con.setRequestProperty("User-Agent", "app");
                con.setRequestProperty("Accept-Language", "en-US,en;q=0.5");
                in = new BufferedReader(new InputStreamReader(con.getInputStream()));
                String inputLine;
                StringBuffer response = new StringBuffer();
                while ((inputLine = in.readLine()) != null) {
                    response.append(inputLine);
                }

                Gson gson = new Gson();
                data = gson.fromJson(response.toString(), UserScore[].class); 
            } catch (Exception e) {
                e.printStackTrace();
            } finally {
                if (in != null) {
                    try {
                        in.close();
                    } catch (IOException e) {
                        e.printStackTrace();
                    }
                }

            }

            return data;
        }

        @Override
        protected void onPostExecute(Object o) {
            super.onPostExecute(o);
            // ...
        }
    }
   ```
   
