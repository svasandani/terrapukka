package api

import (
    "net/http"
    "encoding/json"
    "io/ioutil"

    "github.com/svasandani/terrapukka/internal/util"
    "github.com/svasandani/terrapukka/internal/db"
)

/**************************************************************

API HANDLERS

**************************************************************/

// Handler - export api handler middleware to main
func Handler(w http.ResponseWriter, r *http.Request) {
  // @TODO #4 how do we differentiate between Registrations and
  if r.Method == http.MethodPost {
    if r.URL.Path == "/api/register" {
      createUserHandler(w, r)
    } else if r.URL.Path == "/api/auth" {
      authorizeUserHander(w, r)
    } else {
      util.RespondError(w, 404, "Unknown endpoint: " + r.URL.Path)
    }
  } else {
    util.RespondError(w, 403, "Please use POST requests only.")
  }
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  util.CheckError("Error reading response body:", err)

  user := db.User {}
  err = json.Unmarshal(body, &user)

  util.CheckError("Error unmarshalling response JSON:", err)

  // @TODO #1 get authorization token ?
  db.RegisterUser(user)

  // @TODO #1 write authorization token
  // w.Write(tokenJSON)
}

func authorizeUserHander(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  util.CheckError("Error reading response body:", err)

  auth := db.AuthorizationRequest {}
  err = json.Unmarshal(body, &auth)

  util.CheckError("Error unmarshalling response JSON:", err)

  // @TODO #1 get authorization token ?
  token := db.AuthorizeUser(auth)

  // @TODO #1 write authorization token
  // w.Write(tokenJSON)

  if token.Authorized {
    json, err := json.Marshal(token)

    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    w.Write(json)
  } else {
    util.RespondError(w, 403, "User not found.")
  }
}
