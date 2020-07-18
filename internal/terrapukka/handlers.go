package api

import (
    "net/http"
    "encoding/json"
    "io/ioutil"
)

type httpError struct {
  Status int `json:"status"`
  Msg string `json:"message"`
}

type httpOK struct {
  Msg string `json:"message"`
}

/**************************************************************

HANDLERS

**************************************************************/

// Handler - export handler middleware to main
func Handler(w http.ResponseWriter, r *http.Request) {
  // @TODO #4 how do we differentiate between Registrations and
  if r.Method == http.MethodPost {
    if r.URL.Path == "/register" {
      createUserHandler(w, r)
    } else if r.URL.Path == "/auth" {
      authorizeUserHander(w, r)
    } else {
      respondError(w, 404, "Unknown endpoint: " + r.URL.Path)
    }
  } else {
    respondError(w, 403, "Please use POST requests only.")
  }
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  checkError("Error reading response body:", err)

  user := User {}
  err = json.Unmarshal(body, &user)

  checkError("Error unmarshalling response JSON:", err)

  // @TODO #1 get authorization token ?
  registerUser(user)

  // @TODO #1 write authorization token
  // w.Write(tokenJSON)
}

func authorizeUserHander(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  checkError("Error reading response body:", err)

  auth := AuthorizationRequest {}
  err = json.Unmarshal(body, &auth)

  checkError("Error unmarshalling response JSON:", err)

  // @TODO #1 get authorization token ?
  token := authorizeUser(auth)

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
    respondError(w, 403, "User not found.")
  }
}

// Boilerplate HTTP responses for errors and OKs

func respondError(w http.ResponseWriter, status int, msg string) {
  resp := httpError { Status: status, Msg: msg }

  json, err := json.Marshal(resp)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)

  w.Write(json)
}

func respondOK(w http.ResponseWriter) {
  ok := httpOK { Msg: "Ok" }

  json, err := json.Marshal(ok)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)

  w.Write(json)
}
