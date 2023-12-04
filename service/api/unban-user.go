package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")

	var message string
	pathUserID := ps.ByName("userid")
	if pathUserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User is not correctely authenticated"
		json.NewEncoder(w).Encode(message)
	}

	// find the user and the one to remove
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	var loggedUser User
	var ind int
	for _, user := range Users {
		if user.userID == userID {
			loggedUser = user
			break
		}
	}
	for i, ban := range loggedUser.profile.banned {
		if ban.username == username {
			ind = i
			break
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	loggedUser.profile.banned = append(loggedUser.profile.banned[:ind], loggedUser.profile.banned[ind+1:]...)
	json.NewEncoder(w).Encode(loggedUser)
}
