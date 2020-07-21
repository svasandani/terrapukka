package db

// User - export user struct for http
type User struct {
  Name string `json:"name"`
  Email string `json:"email"`
  Password string `json:"password,omitempty"`
}

// Client - application requesting user data
type Client struct {
  Name string `json:"name,omitempty"`
  ID string `json:"id,omitempty"`
  Secret string `json:"secret,omitempty"`
  RedirectURI string `json:"redirect_uri"`
}

// UserAuthorizationRequest - struct for authenticating users from client
type UserAuthorizationRequest struct {
  ResponseType string `json:"response_type"`
  ClientID string `json:"client_id"`
  RedirectURI string `json:"redirect_uri"`
  State string `json:"state"`
  User User `json:"user"`
}

// UserAuthorizationResponse - struct for responding to user authentication requests
type UserAuthorizationResponse struct {
  RedirectURI string `json:"redirect_uri"`
  AuthCode string `json:"auth_code"`
  State string `json:"state"`
}

// ClientAccessRequest - struct for clients requesting user data
type ClientAccessRequest struct {
  GrantType string `json:"grant_type"`
  AuthCode string `json:"auth_code"`
  Client Client `json:"client"`
}

// ClientAccessResponse - struct for responding to client access request
type ClientAccessResponse struct {
  User User `json:"user"`
}

// ClientIdentificationRequest - struct for requesting client name
type ClientIdentificationRequest struct {
  Client Client `json:"client"`
}

// ClientIdentificationResponse - struct for returning client name
type ClientIdentificationResponse struct {
  Client Client `json:"client"`
}

// Connection - export DBConnection to connect to database
type Connection struct {
  User string
  Password string
  Database string
}
