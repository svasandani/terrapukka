package db

import (
  "fmt"
  "regexp"
  "errors"

  "database/sql"
  "github.com/go-sql-driver/mysql"

  "github.com/svasandani/terrapukka/internal/util"
)

/**************************************************************

DATABASE

**************************************************************/

var db *sql.DB

const ers string = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$"
var er *regexp.Regexp = regexp.MustCompile(ers)

// ConnectDB - connect to the database.
func ConnectDB(dbConn Connection) (*sql.DB) {
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


  ins, err := db.Prepare("INSERT INTO users ( name, email, password, auth_code ) VALUES ( ?, ?, ?, ? )")

  util.CheckError("Error preparing db statement:", err)

  _, err = ins.Exec(user.Name, user.Email, user.Password, "")

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

  err = sel.QueryRow(user.Email, user.Password).Scan(&id, &user.Name, &user.Email, &user.Password, &user.AuthCode)

  util.CheckError("Error executing SELECT statement:", err)

  user.AuthCode = generateAuthCode(user)

  ins, err := db.Prepare("UPDATE users SET auth_code=? WHERE id=?")

  util.CheckError("Error preparing db statement:", err)

  _, err = ins.Exec(user.AuthCode, id)

  util.CheckError("Error executing INSERT statement:", err)

  if id != 0 {
    return AuthorizationToken { Authorized: true, Token: user.AuthCode }, nil
  }

  return AuthorizationToken { Authorized: false, Token: "" }, nil
}

func generateAuthCode(user User) (string) {
  return "0y98hc93hh8hc38h"
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

// RegisterClient - Register the client into the database.
func RegisterClient(client Client) (Client, error) {
  // validate user
  if client.Name == "" {
    return Client{}, errors.New("required field missing: name")
  }

  ins, err := db.Prepare("INSERT INTO clients ( name, identifier, secret ) VALUES ( ?, ?, ? )")

  util.CheckError("Error preparing db statement:", err)

  identifier := "n902hc08yd8014"
  secret := "030f0chhh403h"
  _, err = ins.Exec(client.Name, identifier, secret)

  util.CheckError("Error executing INSERT statement:", err)

  if err != nil {
    return Client { }, errors.New("a problem occurred; please try again later")
  }

  return Client { Name: client.Name, ID: identifier, Secret: secret }, nil
}

// AuthorizeClient - Authorize the client given a ClientAccessRequest
func AuthorizeClient(car ClientAccessRequest) (User, error) {
  if car.AuthCode == "" {
    return User{}, errors.New("required field missing: auth_code")
  }

  if car.Client.ID == "" {
    return User{}, errors.New("required field missing: id")
  }
  if car.Client.Secret == "" {
    return User{}, errors.New("required field missing: secret")
  }

  sel, err := db.Prepare("SELECT * FROM clients WHERE identifier LIKE ? AND secret LIKE ?")
  defer sel.Close()

  var id int
  var client Client

  util.CheckError("Error preparing db statement:", err)

  err = sel.QueryRow(car.Client.ID, car.Client.Secret).Scan(&id, &client.Name, &client.ID, &client.Secret)

  util.CheckError("Error executing SELECT from clients statement:", err)

  if id == 0 {
    return User{}, errors.New("no such client exists")
  }

  var userid int
  var user User

  sel, err = db.Prepare("SELECT * FROM users WHERE auth_code LIKE ?")

  util.CheckError("Error preparing db statement:", err)

  err = sel.QueryRow(car.AuthCode).Scan(&userid, &user.Name, &user.Email, &user.Password, &user.AuthCode)

  util.CheckError("Error executing SELECT FROM users statement:", err)

  if userid == 0 {
    return User{}, errors.New("no such user exists")
  }

  ins, err := db.Prepare("UPDATE users SET auth_code=? WHERE id=?")

  util.CheckError("Error preparing db statement:", err)

  _, err = ins.Exec("", userid)

  util.CheckError("Error executing INSERT statement:", err)

  return user, err
}
