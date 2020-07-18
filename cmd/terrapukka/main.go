package main

import (
  "log"
  "net/http"
  "encoding/json"
  "io/ioutil"

  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type User struct {
  Name string `json:"name"`
  Email string `json:"email"`
  Password string `json:"password"`
}

type HTTPError struct {
  Status int `json:"status"`
  Msg string `json:"message"`
}

type HTTPOk struct {
  Msg string `json:"message"`
}

// @TODO #1 create Authorization struct for OAuth
// type AuthorizationToken struct

var DB *sql.DB

func main() {
  connectDB()
  defer DB.Close()

  http.HandleFunc("/", Handler)

  log.Fatal(http.ListenAndServe(":3000", nil))
}

/**************************************************************

HANDLERS

**************************************************************/

func Handler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
    CreateUserHandler(w, r)
  } else {
    respondError(w, 403, "Please use POST requests only.")
  }
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  checkError("Error reading response body:", err)

  user := User {}
  err = json.Unmarshal(body, &user)

  checkError("Error unmarshalling response JSON:", err)

  // @TODO #1 get authorization token ?
  RegisterUser(user)

  // @TODO #1 write authorization token
  // w.Write(tokenJSON)
}

// Boilerplate HTTP responses for errors and OKs

func respondError(w http.ResponseWriter, status int, msg string) {
  resp := HTTPError { Status: status, Msg: msg }

  json, err := json.Marshal(resp)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)

  w.Write(json)
}

func respondOK(w http.ResponseWriter) {
  ok := HTTPOk { Msg: "Ok" }

  json, err := json.Marshal(ok)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)

  w.Write(json)
}

/**************************************************************

DATABASE

**************************************************************/

// Connect to the database.
func connectDB() {
  // @TODO #2 secrets?
  db, err := sql.Open("mysql", "terrapukka:terrapukka@/terrapukka")

  checkError("Error opening connection to database:", err)

  err = db.Ping()

  checkError("Error establishing connection to database:", err)

  DB = db

  // @QOL create table if not exists, maybe?
}

// Register the user into the database.
func RegisterUser(user User) {
  ct, err := DB.Prepare("INSERT INTO users ( name, email, password ) VALUES ( ?, ?, ? )")

  checkError("Error preparing db statement:", err)

  _, err = ct.Exec(user.Name, user.Email, user.Password)

  checkError("Error executing INSERT statement:", err)
}

/**************************************************************

UTILITY

**************************************************************/

// Check for error; if not nil, print a message along with the error.
func checkError(msg string, err error) {
  if err != nil {
    log.Println(msg)
    log.Println(err.Error())
  }
}
