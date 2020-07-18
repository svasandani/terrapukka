package main

import (
  "log"
  "net/http"

  "github.com/svasandani/terrapukka/internal/api"
  "github.com/svasandani/terrapukka/internal/db"
  "github.com/svasandani/terrapukka/internal/web"
)

func main() {
  database := db.ConnectDB(db.DBConnection { User: "terrapukka", Password: "terrapukka", Database: "terrapukka" })
  defer database.Close()

  web.Init()

  http.HandleFunc("/api/", api.Handler)
  http.HandleFunc("/", web.Handler)

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

  log.Fatal(http.ListenAndServe(":3000", nil))
}
