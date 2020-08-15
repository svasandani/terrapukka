package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type httpError struct {
	Status int    `json:"status"`
	Msg    string `json:"message"`
}

type httpOK struct {
	Msg string `json:"message"`
}

/**************************************************************

UTILITY

**************************************************************/

// CheckError - Check for error; if not nil, print a message along with the error.
func CheckError(msg string, err error) {
	if err != nil {
		log.Println(msg)
		log.Println(err.Error())
	}
}

// RespondError - Boilerplate HTTP responses for error
func RespondError(w http.ResponseWriter, status int, msg string) {
	resp := httpError{Status: status, Msg: msg}

	json, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(json)
}

// RespondOK - Boilerplate HTTP responses for ok
func RespondOK(w http.ResponseWriter) {
	ok := httpOK{Msg: "Ok"}

	json, err := json.Marshal(ok)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(json)
}
