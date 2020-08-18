package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/svasandani/terrapukka/internal/db"
	"github.com/svasandani/terrapukka/internal/util"
)

/**************************************************************

API MIDDLEWARE HANDLERS

**************************************************************/

// Middleware - chain all middleware handlers in one nice convenient function :))
func Middleware(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return CorsHandler(PreflightRequestHandler(PostHandler(JSONHandler(fn))))
}

// CorsHandler - set all CORS headers
func CorsHandler(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		fn(w, r)
	})
}

// PostHandler - ensure all requests to API are posts
func PostHandler(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fn(w, r)
		} else {
			util.RespondError(w, 403, "Please use POST requests only.")
		}
	})
}

// JSONHandler - ensure all requests have JSON payloads
func JSONHandler(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "application/json" {
			fn(w, r)
		} else {
			util.RespondError(w, 400, "Please submit JSON payloads only.")
		}
	})
}

// PreflightRequestHandler - respond with OK on CORS preflight check
func PreflightRequestHandler(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			util.RespondOK(w)
		} else {
			fn(w, r)
		}
	})
}

/**************************************************************

API HANDLERS

**************************************************************/

// CreateUserHandler - create a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	util.CheckError("Error reading response body:", err)

	uar := db.UserAuthorizationRequest{}
	err = json.Unmarshal(body, &uar)

	util.CheckError("Error unmarshalling response JSON:", err)

	resp, err := db.RegisterUser(uar)

	util.CheckError("Error registering user:", err)

	if !(resp == db.UserAuthorizationResponse{}) {
		json, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(json)
	} else {
		util.RespondError(w, 400, err.Error())
	}
}

// AuthorizeUserHandler - authorize a user given email and password
func AuthorizeUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	util.CheckError("Error reading response body:", err)

	uar := db.UserAuthorizationRequest{}
	err = json.Unmarshal(body, &uar)

	util.CheckError("Error unmarshalling response JSON:", err)

	resp, err := db.AuthorizeUser(uar)

	util.CheckError("Error authorizing user:", err)

	if !(resp == db.UserAuthorizationResponse{}) {
		json, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(json)
	} else {
		util.RespondError(w, 403, err.Error())
	}
}

// ResetTokenHandler - handle requests for a reset password token
func ResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	util.CheckError("Error reading response body:", err)

	urtr := db.UserResetTokenRequest{}
	err = json.Unmarshal(body, &urtr)

	util.CheckError("Error unmarshalling response JSON:", err)

	resp, err := db.ResetTokenUser(urtr)

	util.CheckError("Error creating reset password token:", err)

	// send resp via smtp

	if err != nil {
		util.RespondError(w, 400, err.Error())
	} else {
		json, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			w.Write(json)
		}
	}
}

// ResetHandler - handle reset password requests
func ResetHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	util.CheckError("Error reading response body:", err)

	urr := db.UserResetRequest{}
	err = json.Unmarshal(body, &urr)

	util.CheckError("Error unmarshalling response JSON:", err)

	err = db.ResetUser(urr)

	util.CheckError("Error resetting user password:", err)

	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
	} else {
		util.RespondOK(w)
	}
}

// CreateClientHandler - create a new client
func CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	util.CheckError("Error reading response body:", err)

	client := db.Client{}
	err = json.Unmarshal(body, &client)

	util.CheckError("Error unmarshalling response JSON:", err)

	resp, err := db.RegisterClient(client)

	util.CheckError("Error authorizing user:", err)

	if client.ID != "" {
		json, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(json)
	} else {
		util.RespondError(w, 400, err.Error())
	}
}

// AuthorizeClientHandler - authorize a client for a user given ClientAccessRequest
func AuthorizeClientHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	util.CheckError("Error reading response body:", err)

	car := db.ClientAccessRequest{}
	err = json.Unmarshal(body, &car)

	util.CheckError("Error unmarshalling response JSON:", err)

	resp, err := db.AuthorizeClient(car)

	util.CheckError("Error authorizing client:", err)

	if resp.User.Name != "" && resp.User.Email != "" {
		json, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(json)
	} else {
		util.RespondError(w, 403, err.Error())
	}
}

// IdentifyClientHandler - identify a client given ClientIdentificationRequest
func IdentifyClientHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	util.CheckError("Error reading response body:", err)

	cir := db.ClientIdentificationRequest{}
	err = json.Unmarshal(body, &cir)

	util.CheckError("Error unmarshalling response JSON:", err)

	resp, err := db.IdentifyClient(cir)

	util.CheckError("Error identifying client:", err)

	if resp.Client.ID == cir.Client.ID {
		json, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(json)
	} else {
		util.RespondError(w, 403, err.Error())
	}
}
