package db

import (
  "fmt"
  "regexp"
  "errors"

  "database/sql"
  "github.com/go-sql-driver/mysql" // import mysql driver

  "github.com/svasandani/terrapukka/internal/util"
)

/**************************************************************

DATABASE

**************************************************************/

var db *sql.DB

const ers string = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$"
var er *regexp.Regexp = regexp.MustCompile(ers)

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
func RegisterUser(user User) (AuthorizationToken, error) {
  // validate user
  if user.Name == "" {
    return AuthorizationToken{}, errors.New("Required field missing: Name")
  }

  if err := validateEmailPassword(user); err != nil {
    return AuthorizationToken{}, err
  }


  ins, err := db.Prepare("INSERT INTO users ( name, email, password ) VALUES ( ?, ?, ? )")

  util.CheckError("Error preparing db statement:", err)

  _, err = ins.Exec(user.Name, user.Email, user.Password)

  util.CheckError("Error executing INSERT statement:", err)

  if err != nil {
    if sqlErr, ok := err.(*mysql.MySQLError); ok {
      if sqlErr.Number == 1062 {
        return AuthorizationToken { Authorized: false, Token: "" }, errors.New("a user already exists with that email")
      }
    }

    return AuthorizationToken { Authorized: false, Token: "" }, errors.New("a problem occurred; please try again later")
  }

  return AuthorizationToken { Authorized: true, Token: user.Name }, nil
}

// AuthorizeUser - Authorize the user given specific User data
func AuthorizeUser(user User) (AuthorizationToken, error) {
  if err := validateEmailPassword(user); err != nil {
    return AuthorizationToken{}, err
  }

  sel, err := db.Prepare("SELECT * FROM users WHERE email LIKE ? AND password LIKE ?")
  defer sel.Close()

  var id int

  util.CheckError("Error preparing db statement:", err)

  err = sel.QueryRow(user.Email, user.Password).Scan(&id, &user.Name, &user.Email, &user.Password)

  util.CheckError("Error executing SELECT statement:", err)

  if id != 0 {
    return AuthorizationToken { Authorized: true, Token: user.Name }, nil
  }

  return AuthorizationToken { Authorized: false, Token: "" }, nil
}

func validateEmailPassword(user User) (error) {
  if user.Email == "" {
    return errors.New("required field missing: email")
  }
  if !er.MatchString(user.Email) {
    return errors.New("invalid field: email")
  }

  if user.Password == "" {
    return errors.New("required field missing: password")
  }
  if len(user.Password) < 8 {
    return errors.New("password too short; minimum 8 alphanumeric characters")
  }

  return nil
}
