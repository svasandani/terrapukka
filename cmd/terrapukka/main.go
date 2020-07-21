package main

import (
  "log"
  "net/http"
  "flag"

  "github.com/svasandani/terrapukka/internal/api"
  "github.com/svasandani/terrapukka/internal/db"
)

func main() {
  dbuser := flag.String("dbUser", "terrapukka", "Username for MySQL")
  dbpass := flag.String("dbPass", "terrapukka", "Password for MySQL")
  dbname := flag.String("dbName", "terrapukka", "Name of MySQL database")

  port := flag.String("port", "3000", "Port to serve Terrapukka")

  flag.Parse()

  database := db.ConnectDB(db.Connection { User: *dbuser, Password: *dbpass, Database: *dbname })
  defer database.Close()

  http.HandleFunc("/api/register", api.Middleware(api.CreateUserHandler))
  http.HandleFunc("/api/auth", api.Middleware(api.AuthorizeUserHandler))

  http.HandleFunc("/api/client/register", api.Middleware(api.CreateClientHandler))
  http.HandleFunc("/api/client/auth", api.Middleware(api.AuthorizeClientHandler))

  log.Fatal(http.ListenAndServe(":" + *port, nil))
}
