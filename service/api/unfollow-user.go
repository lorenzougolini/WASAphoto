package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	if unfollowedId, ok := UsernameToId[username]; ok {

		getUser := Users[userID]
		if ok2 := contains(getUser.Profile.following, username); ok2 {
			// remove folllowing
			getUser.Profile.following = remove(getUser.Profile.following, username)

			// update followers of the unfollowed user
			unfollowedUser := Users[unfollowedId]
			unfollowedUser.Profile.followers = remove(unfollowedUser.Profile.followers, userID)
		} else {
			w.WriteHeader(http.StatusNotFound)
			message = fmt.Sprintf("The user '%s' is not followed", username)
			json.NewEncoder(w).Encode(message)
			return
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The user '%s' doesn't exist", username)
		json.NewEncoder(w).Encode(message)
		return
	}

	json.NewEncoder(w).Encode(Users[UsernameToId[username]])
}
