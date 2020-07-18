package db

// User - export user struct for http
type User struct {
  Name string `json:"name"`
  Email string `json:"email"`
  Password string `json:"password"`
}

// AuthorizationRequest - export AuthorizationRequest struct for http
type AuthorizationRequest struct {
  Email string `json:"email"`
  Password string `json:"password"`
}

// @TODO #1 create Authorization struct for OAuth
// type AuthorizationToken struct
type AuthorizationToken struct {
  Authorized bool `json:"authorized"`
  Token string `json:"token"`
}

// DBConnection - export DBConnection to connect to database
type DBConnection struct {
  User string
  Password string
  Database string
}
