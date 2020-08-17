package db

import (
	"errors"
	"fmt"
	"regexp"
	"time"

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
func ConnectDB(dbConn Connection) *sql.DB {
	conn := fmt.Sprintf("%v:%v@/%v?parseTime=true", dbConn.User, dbConn.Password, dbConn.Database)

	dbLocal, err := sql.Open("mysql", conn)

	util.CheckError("Error opening connection to database:", err)

	err = dbLocal.Ping()

	util.CheckError("Error establishing connection to database:", err)

	db = dbLocal

	return dbLocal

	// @QOL create table if not exists, maybe?
}

// RegisterUser - Register the user into the database.
func RegisterUser(uar UserAuthorizationRequest) (UserAuthorizationResponse, error) {
	// validate request
	if uar.User.Name == "" {
		return UserAuthorizationResponse{}, errors.New("required field missing: name")
	}

	if err := validateUserAuthorizationRequest(uar); err != nil {
		return UserAuthorizationResponse{}, err
	}

	sel, err := db.Prepare("SELECT id FROM clients WHERE identifier LIKE ? AND redirect_uri LIKE ?")
	defer sel.Close()

	var id int

	util.CheckError("Error preparing db statement:", err)

	err = sel.QueryRow(uar.ClientID, uar.RedirectURI).Scan(&id)

	util.CheckError("Error executing SELECT from clients statement:", err)

	if id == 0 {
		return UserAuthorizationResponse{}, errors.New("client could not be found")
	}

	authCode := generateAuthCode(uar.User)
	password, err := util.HashAndSalt(uar.User.Password)

	util.CheckError("Error salting user password:", err)

	if err != nil {
		return UserAuthorizationResponse{}, errors.New("a problem occurred; please try again later")
	}

	ins, err := db.Prepare("INSERT INTO users ( name, email, password, auth_code, auth_code_generated_at ) VALUES ( ?, ?, ?, ?, ? )")

	util.CheckError("Error preparing db statement:", err)

	_, err = ins.Exec(uar.User.Name, uar.User.Email, password, authCode, time.Now())

	util.CheckError("Error executing INSERT statement:", err)

	if err != nil {
		if sqlErr, ok := err.(*mysql.MySQLError); ok {
			if sqlErr.Number == 1062 {
				return UserAuthorizationResponse{}, errors.New("a user already exists with that email")
			}
		}

		return UserAuthorizationResponse{}, errors.New("a problem occurred; please try again later")
	}

	return UserAuthorizationResponse{RedirectURI: uar.RedirectURI, AuthCode: authCode, State: uar.State}, nil
}

// AuthorizeUser - Authorize the user given specific User data
func AuthorizeUser(uar UserAuthorizationRequest) (UserAuthorizationResponse, error) {
	if uar.ResponseType == "code" {
		return codeAuthorizeUser(uar)
	} else {
		return UserAuthorizationResponse{}, fmt.Errorf("unknown response type: %v", uar.ResponseType)
	}
}

func codeAuthorizeUser(uar UserAuthorizationRequest) (UserAuthorizationResponse, error) {
	// validate request
	if err := validateUserAuthorizationRequest(uar); err != nil {
		return UserAuthorizationResponse{}, err
	}

	sel, err := db.Prepare("SELECT id FROM clients WHERE identifier LIKE ? AND redirect_uri LIKE ?")
	defer sel.Close()

	var id int

	util.CheckError("Error preparing db statement:", err)

	err = sel.QueryRow(uar.ClientID, uar.RedirectURI).Scan(&id)

	util.CheckError("Error executing SELECT from clients statement:", err)

	if id == 0 {
		return UserAuthorizationResponse{}, errors.New("client could not be found")
	}

	sel, err = db.Prepare("SELECT id, password FROM users WHERE email LIKE ?")
	defer sel.Close()

	var userid int
	var password string
	var authCode string

	util.CheckError("Error preparing db statement:", err)

	err = sel.QueryRow(uar.User.Email).Scan(&userid, &password)

	util.CheckError("Error executing SELECT statement:", err)

	err = util.CompareHashAndText(password, uar.User.Password)

	if err != nil {
		return UserAuthorizationResponse{}, errors.New("user email or password is incorrect")
	}

	authCode = generateAuthCode(uar.User)

	ins, err := db.Prepare("UPDATE users SET auth_code=?, auth_code_generated_at=? WHERE id=?")

	util.CheckError("Error preparing db statement:", err)

	_, err = ins.Exec(authCode, time.Now(), userid)

	util.CheckError("Error executing INSERT statement:", err)

	if userid != 0 {
		return UserAuthorizationResponse{RedirectURI: uar.RedirectURI, AuthCode: authCode, State: uar.State}, nil
	}

	return UserAuthorizationResponse{}, errors.New("user could not be found")
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

	identifier := util.UUID()
	secret := util.UUID()
	hashed, err := util.HashAndSalt(secret)

	util.CheckError("Error salting client secret:", err)

	if err != nil {
		return Client{}, errors.New("a problem occurred; please try again later")
	}

	_, err = ins.Exec(client.Name, identifier, hashed, client.RedirectURI)

	util.CheckError("Error executing INSERT statement:", err)

	if err != nil {
		return Client{}, errors.New("a problem occurred; please try again later")
	}

	return Client{Name: client.Name, ID: identifier, Secret: secret, RedirectURI: client.RedirectURI}, nil
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

	if car.GrantType == "identity" {
		return identityAuthorizeClient(car)
	} else {
		return ClientAccessResponse{}, fmt.Errorf("unknown grant type: %v", car.GrantType)
	}
}

func identityAuthorizeClient(car ClientAccessRequest) (ClientAccessResponse, error) {
	sel, err := db.Prepare("SELECT id, secret FROM clients WHERE identifier LIKE ? AND redirect_uri LIKE ?")
	defer sel.Close()

	var id int
	var secret string

	util.CheckError("Error preparing db statement:", err)

	err = sel.QueryRow(car.Client.ID, car.Client.RedirectURI).Scan(&id, &secret)

	util.CheckError("Error executing SELECT from clients statement:", err)

	if id == 0 {
		return ClientAccessResponse{}, errors.New("client could not be found")
	}

	err = util.CompareHashAndText(secret, car.Client.Secret)

	if err != nil {
		return ClientAccessResponse{}, errors.New("client secret is incorrect")
	}

	var userid int
	var timeset time.Time
	var user User

	sel, err = db.Prepare("SELECT id, name, email, auth_code_generated_at FROM users WHERE auth_code LIKE ?")

	util.CheckError("Error preparing user select db statement:", err)

	err = sel.QueryRow(car.AuthCode).Scan(&userid, &user.Name, &user.Email, &timeset)

	util.CheckError("Error executing SELECT FROM users statement:", err)

	timeset = timeset.Add(15 * time.Minute)

	if time.Now().After(timeset) {
		return ClientAccessResponse{}, errors.New("auth code expired, please sign in again")
	}

	if userid == 0 {
		return ClientAccessResponse{}, errors.New("user could not be found")
	}

	sel, err = db.Prepare("SELECT * FROM roles WHERE user_id=?")

	util.CheckError("Error preparing role select db statement:", err)

	rows, err := sel.Query(userid)

	util.CheckError("Error executing SELECT FROM roles statement:", err)

	defer rows.Close()

	for rows.Next() {
		var role Role
		var roleusr int

		err = rows.Scan(&roleusr, &role.ResourceID, &role.Type)

		util.CheckError("Error scanning result into Role:", err)

		if roleusr != userid {
			return ClientAccessResponse{}, errors.New("user id is not user id ???")
		}

		user.Roles = append(user.Roles, role)
	}

	ins, err := db.Prepare("UPDATE users SET auth_code=?, auth_code_generated_at=NULL WHERE id=?")

	util.CheckError("Error preparing db statement:", err)

	_, err = ins.Exec("", userid)

	util.CheckError("Error executing INSERT statement:", err)

	return ClientAccessResponse{User: user}, err
}

// IdentifyClient - Identify a client given a ClientAccessRequest
func IdentifyClient(cir ClientIdentificationRequest) (ClientIdentificationResponse, error) {
	if cir.Client.ID == "" {
		return ClientIdentificationResponse{}, errors.New("required field missing: id")
	}
	if cir.Client.RedirectURI == "" {
		return ClientIdentificationResponse{}, errors.New("required field missing: redirect_uri")
	}

	sel, err := db.Prepare("SELECT id, identifier, redirect_uri, name FROM clients WHERE identifier LIKE ? AND redirect_uri LIKE ?")
	defer sel.Close()

	var id int
	var client Client

	util.CheckError("Error preparing db statement:", err)

	err = sel.QueryRow(cir.Client.ID, cir.Client.RedirectURI).Scan(&id, &client.ID, &client.RedirectURI, &client.Name)

	util.CheckError("Error executing SELECT from clients statement:", err)

	if id == 0 {
		return ClientIdentificationResponse{}, errors.New("client could not be found")
	}

	return ClientIdentificationResponse{Client: client}, err
}

func generateAuthCode(user User) string {
	return util.UUID()
}

/**************************************************************

VALIDATOR FUNCTIONS

**************************************************************/

func validateEmailPassword(user User) error {
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

func validateUserAuthorizationRequest(uar UserAuthorizationRequest) error {
	if err := validateEmailPassword(uar.User); err != nil {
		return err
	}

	if uar.ResponseType == "" {
		return errors.New("required field missing: response_type")
	}
	if uar.ClientID == "" {
		return errors.New("required field missing: client_id")
	}
	if uar.RedirectURI == "" {
		return errors.New("required field missing: redirect_uri")
	}

	return nil
}
