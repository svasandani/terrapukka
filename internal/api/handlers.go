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
  return PreflightRequestHandler(PostHandler(JSONHandler(fn)))
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

// PreflightRequestHandler - respond with OK on CORS preflight check
func PreflightRequestHandler(fn func(w http.ResponseWriter, r *http.Request)) (func(w http.ResponseWriter, r *http.Request)) {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
      util.RespondOK(w)
    } else {
      fn(w, r)
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

  token, err := db.RegisterUser(user)

  util.CheckError("Error authorizing user:", err)

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

  uar := db.UserAuthenticationRequest{}
  err = json.Unmarshal(body, &uar)

  util.CheckError("Error unmarshalling response JSON:", err)

  resp, err := db.AuthorizeUser(uar)

  util.CheckError("Error authorizing user:", err)

  if !(resp == db.UserAuthenticationResponse{}) {
    json, err := json.Marshal(resp)

    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    w.Write(json)
  } else {
    util.RespondError(w, 403, "user not found")
  }
}

// CreateClientHandler - create a new client
func CreateClientHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  util.CheckError("Error reading response body:", err)

  client := db.Client {}
  err = json.Unmarshal(body, &client)

  util.CheckError("Error unmarshalling response JSON:", err)

  client, err = db.RegisterClient(client)

  util.CheckError("Error authorizing user:", err)

  if client.ID != "" {
    json, err := json.Marshal(client)

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

// AuthorizeClientHandler - authorize a client for a user given ClientAccessRequest
func AuthorizeClientHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  util.CheckError("Error reading response body:", err)

  car := db.ClientAccessRequest {}
  err = json.Unmarshal(body, &car)

  util.CheckError("Error unmarshalling response JSON:", err)

  resp, err := db.AuthorizeClient(car)

  util.CheckError("Error authorizing client:", err)

  if (resp.User != db.User{}) {
    json, err := json.Marshal(resp)

    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    w.Write(json)
  } else {
    util.RespondError(w, 403, err.Error())
  }
}
