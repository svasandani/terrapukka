package web

import (
  "fmt"
  "net/http"
)

/**************************************************************

WEB HANDLERS

**************************************************************/

// Handler - export web handler middleware to main
func Handler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("HELLO")
}
