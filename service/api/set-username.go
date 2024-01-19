package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")

	var message string
	// check Bearer token
	if !checkLogin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	newUsername, err := io.ReadAll(r.Body)
	if err != nil || string(newUsername) == "" || len(string(newUsername)) < 3 || len(string(newUsername)) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", newUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	// // retrieve user and change its username
	// getUser := Users[userID]

	// UsernameToId[string(newUsername)] = userID // change and delete username->id map
	// delete(UsernameToId, getUser.Username)

	// getUser.Username = string(newUsername) // change username in Users map
	// Users[userID] = getUser

	// Logged["logged"] = getUser.UserID // change username in logged data

	user, err := rt.db.GetByUsername(Logged.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = "Error changing the username"
		_ = json.NewEncoder(w).Encode(message)
		return
	} else {
		err := rt.db.SetName(user.UserID, string(newUsername))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
		Logged.Username = string(newUsername)
	}
	_ = json.NewEncoder(w).Encode(Logged)
}
