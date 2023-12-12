package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// if followedId, ok := UsernameToId[username]; ok {

	// 	// add folllowing
	// 	getUser := Users[userID]
	// 	getUser.Profile.following = append(getUser.Profile.following, username)

	// 	// update followers of the folllowed user
	// 	followedUser := Users[followedId]
	// 	followedUser.Profile.followers = append(followedUser.Profile.followers, userID)

	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	message = fmt.Sprintf("The username '%s' doesn't exist", username)
	// 	json.NewEncoder(w).Encode(message)
	// 	return
	// }

	_, err := rt.db.GetUser(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The username '%s' doesn't exist", username)
		json.NewEncoder(w).Encode(message)
		return
	} else {
		err := rt.db.FollowUser(Logged, username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		message = Logged["username"] + " succesfully followed: " + username
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}

	json.NewEncoder(w).Encode(Users[UsernameToId[username]])
}
