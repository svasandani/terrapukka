package api

import (
  "net/http"
  "encoding/json"
  "io/ioutil"

  "github.com/svasandani/terrapukka/internal/util"
  "github.com/svasandani/terrapukka/internal/db"
)

/**************************************************************

API MIDDLEWARE HANDLERS

**************************************************************/

// Middleware - chain all middleware handlers in one nice convenient function :))
func Middleware(fn func(w http.ResponseWriter, r *http.Request)) (func(w http.ResponseWriter, r *http.Request)) {
  return PostHandler(JSONHandler(fn))
}

// PostHandler - ensure all requests to API are posts
func PostHandler(fn func(w http.ResponseWriter, r *http.Request)) (func(w http.ResponseWriter, r *http.Request)) {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
      fn(w, r)
    } else {
      util.RespondError(w, 403, "Please use POST requests only.")
    }
  })
}

// JSONHandler - ensure all requests have JSON payloads
func JSONHandler(fn func(w http.ResponseWriter, r *http.Request)) (func(w http.ResponseWriter, r *http.Request)) {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Content-Type") == "application/json" {
      fn(w, r)
    } else {
      util.RespondError(w, 400, "Please submit JSON payloads only.")
    }
  })
}

/**************************************************************

API HANDLERS

**************************************************************/

// CreateUserHandler - create a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  util.CheckError("Error reading response body:", err)

  user := db.User {}
  err = json.Unmarshal(body, &user)

  util.CheckError("Error unmarshalling response JSON:", err)

  // @TODO #1 get authorization token ?
  token, err := db.RegisterUser(user)

  util.CheckError("Error authorizing user:", err)

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
    util.RespondError(w, 400, err.Error())
  }
}

// AuthorizeUserHandler - authorize a user given email and password
func AuthorizeUserHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  util.CheckError("Error reading response body:", err)

  user := db.User {}
  err = json.Unmarshal(body, &user)

  util.CheckError("Error unmarshalling response JSON:", err)

  // @TODO #1 get authorization token ?
  token, err := db.AuthorizeUser(user)

  util.CheckError("Error authorizing user:", err)

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
