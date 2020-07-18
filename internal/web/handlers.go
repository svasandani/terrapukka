package web

import (
  "net/http"

  "github.com/svasandani/terrapukka/internal/util"
)

/**************************************************************

WEB HANDLERS

**************************************************************/

// Handler - export web handler middleware to main
func Handler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    if r.URL.Path == "/sign_up" {
      // signUpHandler(w, r)
    } else if r.URL.Path == "/sign_in" {
      // signInHandler(w, r)
    } else {
      util.RespondError(w, 404, "The directory you're looking for couldn't be found.")
    }
  } else {
    util.RespondError(w, 403, "Please use GET requests only.")
  }
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

  err := tmpl.ExecuteTemplate(w, "sign_in.html", nil)

  util.CheckError("Error executing header template:", err)
}
