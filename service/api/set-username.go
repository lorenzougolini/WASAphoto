package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) setUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")

	var message string
	pathUserID := ps.ByName("userid")
	if pathUserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User is not correctely authenticated"
		json.NewEncoder(w).Encode(message)
	}

	username, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// look for logged user and change its username
	var loggedUser User
	for _, user := range Users {
		if user.userID == uuid.Must(uuid.Parse(pathUserID)) {
			loggedUser = user
			break
		}
	}
	loggedUser.username = string(username)
	json.NewEncoder(w).Encode(loggedUser)
}
