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
	userID := ps.ByName("userid")
	// check logged user id
	if !checkLogin(userID) {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User is not correctly authenticated"
		json.NewEncoder(w).Encode(message)
		return
	} else if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUsername, err := io.ReadAll(r.Body)
	if err != nil || string(newUsername) == "" || len(string(newUsername)) < 3 || len(string(newUsername)) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", newUsername)
		json.NewEncoder(w).Encode(message)
		return
	}

	// retrieve user and change its username
	getUser := Users[userID]

	UsernameToId[string(newUsername)] = userID // change and delete username->id map
	delete(UsernameToId, getUser.Username)

	getUser.Username = string(newUsername) // change username in Users map
	Users[userID] = getUser

	Logged["logged"] = getUser.UserID // change username in logged data

	json.NewEncoder(w).Encode(Logged["logged"])
}