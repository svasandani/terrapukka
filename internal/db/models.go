package db

// User - export user struct for http
type User struct {
  Name string `json:"name,omitempty"`
  Email string `json:"email"`
  Password string `json:"password"`
}

// @TODO #1 create Authorization struct for OAuth
// AuthorizationToken - token for OAuth authorization
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
