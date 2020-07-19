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

DATABASE FUNCTIONS

**************************************************************/

var db *sql.DB

const ers string = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$"
var er *regexp.Regexp = regexp.MustCompile(ers)

// ConnectDB - connect to the database.
func ConnectDB(dbConn Connection) (*sql.DB) {
  conn := fmt.Sprintf("%v:%v@/%v", dbConn.User, dbConn.Password, dbConn.Database)

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
func AuthorizeUser(uar UserAuthenticationRequest) (UserAuthenticationResponse, error) {
  if err := validateEmailPassword(uar.User); err != nil {
    return UserAuthenticationResponse{}, err
  }

  if uar.ClientID == "" {
    return UserAuthenticationResponse{}, errors.New("required field missing: client_id")
  }
  if uar.RedirectURI == "" {
    return UserAuthenticationResponse{}, errors.New("required field missing: redirect_uri")
  }

  sel, err := db.Prepare("SELECT id FROM clients WHERE identifier LIKE ? AND redirect_uri LIKE ?")
  defer sel.Close()

  var id int

  util.CheckError("Error preparing db statement:", err)

  err = sel.QueryRow(uar.ClientID, uar.RedirectURI).Scan(&id)

  util.CheckError("Error executing SELECT from clients statement:", err)

  if id == 0 {
    return UserAuthenticationResponse{}, errors.New("no such client exists")
  }

  sel, err = db.Prepare("SELECT id FROM users WHERE email LIKE ? AND password LIKE ?")
  defer sel.Close()

  var userid int
  var authCode string

  util.CheckError("Error preparing db statement:", err)

  err = sel.QueryRow(uar.User.Email, uar.User.Password).Scan(&userid)

  util.CheckError("Error executing SELECT statement:", err)

  authCode = generateAuthCode(uar.User)

  ins, err := db.Prepare("UPDATE users SET auth_code=? WHERE id=?")

  util.CheckError("Error preparing db statement:", err)

  _, err = ins.Exec(authCode, userid)

  util.CheckError("Error executing INSERT statement:", err)

  if userid != 0 {
    return UserAuthenticationResponse { RedirectURI: uar.RedirectURI, AuthCode: authCode, State: uar.State }, nil
  }

  return UserAuthenticationResponse {}, errors.New("user could not be found")
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
  if client.RedirectURI == "" {
    return Client{}, errors.New("required field missing: redirect_uri")
  }

  ins, err := db.Prepare("INSERT INTO clients ( name, identifier, secret, redirect_uri ) VALUES ( ?, ?, ?, ? )")

  util.CheckError("Error preparing db statement:", err)

  identifier := "n902hc08yd8014"
  secret := "030f0chhh403h"
  _, err = ins.Exec(client.Name, identifier, secret, client.RedirectURI)

  util.CheckError("Error executing INSERT statement:", err)

  if err != nil {
    return Client { }, errors.New("a problem occurred; please try again later")
  }

  return Client { Name: client.Name, ID: identifier, Secret: secret, RedirectURI: client.RedirectURI }, nil
}

// AuthorizeClient - Authorize the client given a ClientAccessRequest
func AuthorizeClient(car ClientAccessRequest) (ClientAccessResponse, error) {
  if car.AuthCode == "" {
    return ClientAccessResponse{}, errors.New("required field missing: auth_code")
  }

  if car.Client.ID == "" {
    return ClientAccessResponse{}, errors.New("required field missing: id")
  }
  if car.Client.Secret == "" {
    return ClientAccessResponse{}, errors.New("required field missing: secret")
  }
  if car.Client.RedirectURI == "" {
    return ClientAccessResponse{}, errors.New("required field missing: redirect_uri")
  }

  sel, err := db.Prepare("SELECT id FROM clients WHERE identifier LIKE ? AND secret LIKE ? AND redirect_uri LIKE ?")
  defer sel.Close()

  var id int

  util.CheckError("Error preparing db statement:", err)

  err = sel.QueryRow(car.Client.ID, car.Client.Secret, car.Client.RedirectURI).Scan(&id)

  util.CheckError("Error executing SELECT from clients statement:", err)

  if id == 0 {
    return ClientAccessResponse{}, errors.New("no such client exists")
  }

  var userid int
  var user User

  sel, err = db.Prepare("SELECT id, name, email FROM users WHERE auth_code LIKE ?")

  util.CheckError("Error preparing db statement:", err)

  err = sel.QueryRow(car.AuthCode).Scan(&userid, &user.Name, &user.Email)

  util.CheckError("Error executing SELECT FROM users statement:", err)

  if userid == 0 {
    return ClientAccessResponse{}, errors.New("no such user exists")
  }

  ins, err := db.Prepare("UPDATE users SET auth_code=? WHERE id=?")

  util.CheckError("Error preparing db statement:", err)

  _, err = ins.Exec("", userid)

  util.CheckError("Error executing INSERT statement:", err)

  return ClientAccessResponse{User: user}, err
}
