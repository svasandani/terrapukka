package db

// User - export user struct for http
type User struct {
  Name string `json:"name,omitempty"`
  Email string `json:"email"`
  Password string `json:"password"`
}

// Client - application requesting user data
type Client struct {
  Name string `json:"name,omitempty"`
  ID string `json:"id,omitempty"`
  Secret string `json:"secret,omitempty"`
  RedirectURI string `json:"redirect_uri"`
}

// ClientAccessRequest - struct for clients requesting user data
type ClientAccessRequest struct {
  GrantType string `json:"grant_type"`
  AuthCode string `json:"auth_code"`
  Client Client `json:"client"`
}

// ClientAccessResponse - stsruct for responding to client access request
type ClientAccessResponse struct {
  User User `json:"user"`
}

// UserAuthenticationRequest - struct for authenticating users from client
type UserAuthenticationRequest struct {
  ResponseType string `json:"response_type"`
  ClientID string `json:"client_id"`
  RedirectURI string `json:"redirect_uri"`
  State string `json:"state"`
  User User `json:"user"`
}

// UserAuthenticationResponse - struct for responding to user authentication requests
type UserAuthenticationResponse struct {
  RedirectURI string `json:"redirect_uri"`
  AuthCode string `json:"auth_code"`
  State string `json:"state"`
}

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
