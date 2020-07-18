package api

import (
  "fmt"

  "database/sql"
  _ "github.com/go-sql-driver/mysql" // import mysql driver
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

  checkError("Error opening connection to database:", err)

  err = dbLocal.Ping()

  checkError("Error establishing connection to database:", err)

  db = dbLocal

  return dbLocal

  // @QOL create table if not exists, maybe?
}

// Register the user into the database.
func registerUser(user User) {
  ins, err := db.Prepare("INSERT INTO users ( name, email, password ) VALUES ( ?, ?, ? )")

  checkError("Error preparing db statement:", err)

  _, err = ins.Exec(user.Name, user.Email, user.Password)

  checkError("Error executing INSERT statement:", err)
}

func authorizeUser(auth AuthorizationRequest) (AuthorizationToken) {
  sel, err := db.Prepare("SELECT * FROM users WHERE email LIKE ? AND password LIKE ?")
  defer sel.Close()

  var user User
  var id int

  checkError("Error preparing db statement:", err)

  err = sel.QueryRow(auth.Email, auth.Password).Scan(&id, &user.Name, &user.Email, &user.Password)

  checkError("Error executing INSERT statement:", err)

  if id != 0 {
    return AuthorizationToken { Authorized: true, Token: "test" }
  }

  return AuthorizationToken { Authorized: false, Token: "" }
}
