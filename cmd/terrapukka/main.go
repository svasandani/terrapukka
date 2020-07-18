package main

import (
  "log"
  "net/http"

  "github.com/svasandani/terrapukka/internal/terrapukka"
)

func main() {
  db := api.ConnectDB(api.DBConnection { User: "terrapukka", Password: "terrapukka", Database: "terrapukka" })
  defer db.Close()

  http.HandleFunc("/", api.Handler)

  log.Fatal(http.ListenAndServe(":3000", nil))
}
