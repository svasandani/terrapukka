package db

import (
  "fmt"

  "database/sql"
  _ "github.com/go-sql-driver/mysql" // import mysql driver

  "github.com/svasandani/terrapukka/internal/util"
)

/**************************************************************

DATABASE

**************************************************************/

var db *sql.DB

// ConnectDB - connect to the database.
func ConnectDB(dbConn DBConnection) (*sql.DB) {
  conn := fmt.Sprintf("%v:%v@/%v", dbConn.User, dbConn.Password, dbConn.Database)

  // @TODO #2 secrets?
  dbLocal, err := sql.Open("mysql", conn)

  util.CheckError("Error opening connection to database:", err)

  err = dbLocal.Ping()

  util.CheckError("Error establishing connection to database:", err)

  db = dbLocal

  return dbLocal

  // @QOL create table if not exists, maybe?
}

// RegisterUser - Register the user into the database.
func RegisterUser(user User) {
  ins, err := db.Prepare("INSERT INTO users ( name, email, password ) VALUES ( ?, ?, ? )")

  util.CheckError("Error preparing db statement:", err)

  _, err = ins.Exec(user.Name, user.Email, user.Password)

  util.CheckError("Error executing INSERT statement:", err)
}

// AuthorizeUser - Authorize the user given a specific AuthorizationRequest
func AuthorizeUser(auth AuthorizationRequest) (AuthorizationToken) {
  sel, err := db.Prepare("SELECT * FROM users WHERE email LIKE ? AND password LIKE ?")
  defer sel.Close()

  var user User
  var id int

  util.CheckError("Error preparing db statement:", err)

  err = sel.QueryRow(auth.Email, auth.Password).Scan(&id, &user.Name, &user.Email, &user.Password)

  util.CheckError("Error executing SELECT statement:", err)

  if id != 0 {
    return AuthorizationToken { Authorized: true, Token: "test" }
  }

  return AuthorizationToken { Authorized: false, Token: "" }
}
