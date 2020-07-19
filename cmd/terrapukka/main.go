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

  http.HandleFunc("/api/register", api.Middleware(api.CreateUserHandler))
  http.HandleFunc("/api/auth", api.Middleware(api.AuthorizeUserHandler))

  http.HandleFunc("/api/client/register", api.Middleware(api.CreateClientHandler))
  // expect AUTH_CODE, CLIENT_ID, CLIENT_SECRET, and return name and email corresponding to AUTH_CODE
  http.HandleFunc("/api/client/auth", api.Middleware(api.AuthorizeClientHandler))

  http.HandleFunc("/", web.Handler)
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

  log.Fatal(http.ListenAndServe(":" + *port, nil))
}
