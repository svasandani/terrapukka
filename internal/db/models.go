package db

// User - export user struct for http
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Roles    []Role `json:"roles,omitempty"`
}

// Client - application requesting user data
type Client struct {
	Name        string `json:"name,omitempty"`
	ID          string `json:"id,omitempty"`
	Secret      string `json:"secret,omitempty"`
	RedirectURI string `json:"redirect_uri"`
}

// Role - export generic Role struct
type Role struct {
	Type       string `json:"type"`
	ResourceID string `json:"resource_id"`
}

// UserAuthorizationRequest - struct for authenticating users from client
type UserAuthorizationRequest struct {
	ResponseType string `json:"response_type"`
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	State        string `json:"state"`
	User         User   `json:"user"`
}

// UserAuthorizationResponse - struct for responding to user authentication requests
type UserAuthorizationResponse struct {
	RedirectURI string `json:"redirect_uri"`
	AuthCode    string `json:"auth_code"`
	State       string `json:"state"`
}

// UserResetTokenRequest - struct for responding to password reset requests
type UserResetTokenRequest struct {
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	User        User   `json:"user"`
}

// UserResetTokenResponse - struct to return password reset token via smtp
type UserResetTokenResponse struct {
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	User        User   `json:"user"`
	ResetToken  string `json:"token"`
}

// UserResetRequest - struct for handling password reset with token
type UserResetRequest struct {
	ResetToken string `json:"reset_token"`
	User       User   `json:"user"`
}

// ClientAccessRequest - struct for clients requesting user data
type ClientAccessRequest struct {
	GrantType string `json:"grant_type"`
	AuthCode  string `json:"auth_code"`
	Client    Client `json:"client"`
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
	User     string
	Password string
	Database string
}
