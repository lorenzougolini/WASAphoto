package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// check username to follow and proceed it exists
	username := ps.ByName("username")
	if (username) == "" || len(string(username)) < 3 || len(string(username)) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", username)
		json.NewEncoder(w).Encode(message)
		return
	}

	// if _, ok := UsernameToId[username]; ok {

	// 	// add ban
	// 	getUser := Users[userID]
	// 	getUser.Profile.banned = append(getUser.Profile.banned, username)

	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	message = fmt.Sprintf("The username '%s' doesn't exist", username)
	// 	json.NewEncoder(w).Encode(message)
	// 	return
	// }

	user, err := rt.db.GetUser(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The username '%s' doesn't exist", username)
		json.NewEncoder(w).Encode(message)
		return
	} else {
		err := rt.db.BanUser(Logged, user.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		message = Logged["username"] + " succesfully banned: " + user.Username
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}

	json.NewEncoder(w).Encode(Users[UsernameToId[username]])
}
