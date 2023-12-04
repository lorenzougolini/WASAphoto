package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")

	var message string
	pathUserID := ps.ByName("userid")
	if pathUserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User is not correctely authenticated"
		json.NewEncoder(w).Encode(message)
	}

	// find user to follow
	username, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if string(username) == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	var loggedUser User
	var banUser User
	for _, user := range Users {
		if user.userID == userID {
			loggedUser = user
		}
	}
	for _, user := range Users {
		if user.username == string(username) {
			banUser = user
			break
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	loggedUser.profile.banned = append(loggedUser.profile.following, banUser)
	json.NewEncoder(w).Encode(loggedUser)
}
