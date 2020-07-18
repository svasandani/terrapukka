package main

import (
  "log"
  "net/http"
  "flag"

  "github.com/svasandani/terrapukka/internal/api"
  "github.com/svasandani/terrapukka/internal/db"
  "github.com/svasandani/terrapukka/internal/web"
)

func main() {
  dbuser := flag.String("dbUser", "terrapukka", "Username for MySQL")
  dbpass := flag.String("dbPass", "terrapukka", "Password for MySQL")
  dbname := flag.String("dbName", "terrapukka", "Name of MySQL database")

  port := flag.String("port", "3000", "Port to serve Terrapukka")

  flag.Parse()

  database := db.ConnectDB(db.Connection { User: *dbuser, Password: *dbpass, Database: *dbname })
  defer database.Close()

  web.Init()

  http.HandleFunc("/api/", api.Handler)
  http.HandleFunc("/", web.Handler)

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

  log.Fatal(http.ListenAndServe(":" + *port, nil))
}
