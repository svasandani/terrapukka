package db

// User - export user struct for http
type User struct {
  Name string `json:"name,omitempty"`
  Email string `json:"email"`
  Password string `json:"password"`
  AuthCode string `json:"auth_code,omitempty"`
}

// Client - application requesting user data
type Client struct {
  Name string `json:"name,omitempty"`
  ID string `json:"id,omitempty"`
  Secret string `json:"secret,omitempty"`
}

// ClientAccessRequest - struct for clients requesting user data
type ClientAccessRequest struct {
  AuthCode string `json:"auth_code"`
  Client Client `json:"client"`
}

// @TODO #1 create Authorization struct for OAuth
// AuthorizationToken - token for OAuth authorization
type AuthorizationToken struct {
  Authorized bool `json:"authorized"`
  Token string `json:"token"`
}

// Connection - export DBConnection to connect to database
type Connection struct {
  User string
  Password string
  Database string
}
