package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var message string
	// check Bearer token
	if !checkLogin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	stream, err := rt.db.GetStream(Logged.UserID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		message = "Error getting your stream"
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	streamJson, err := json.MarshalIndent(stream, "", " ")
	if err != nil {
		message = "Error retrieving the profile"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	// _ = json.NewEncoder(w).Encode(Logged.Username)
	_, _ = w.Write(streamJson)
}
