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
	username := ps.ByName("username")
	// check Bearer token
	if !checkLogin(r) || username != Logged.Username {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	// check username to unfollow and proceed it exists
	unfollowedUsername := ps.ByName("unfollowedUsername")
	if unfollowedUsername == username || unfollowedUsername == "" || len(unfollowedUsername) < 3 || len(unfollowedUsername) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", unfollowedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	// if unfollowedId, ok := UsernameToId[username]; ok {

	// 	getUser := Users[userID]
	// 	if ok2 := contains(getUser.Profile.following, username); ok2 {
	// 		// remove folllowing
	// 		getUser.Profile.following = remove(getUser.Profile.following, username)

	// 		// update followers of the unfollowed user
	// 		unfollowedUser := Users[unfollowedId]
	// 		unfollowedUser.Profile.followers = remove(unfollowedUser.Profile.followers, userID)
	// 	} else {
	// 		w.WriteHeader(http.StatusNotFound)
	// 		message = fmt.Sprintf("The user '%s' is not followed", username)
	// 		json.NewEncoder(w).Encode(message)
	// 		return
	// 	}

	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	message = fmt.Sprintf("The user '%s' doesn't exist", username)
	// 	json.NewEncoder(w).Encode(message)
	// 	return
	// }

	unfollowedUser, err := rt.db.GetByUsername(unfollowedUsername)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The username '%s' doesn't exist", unfollowedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return

	} else {
		err := rt.db.UnfollowUser(Logged.UserID, unfollowedUser.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
		message = Logged.Username + " succesfully unfollowed: " + unfollowedUsername
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(message)
	}

	// _ = json.NewEncoder(w).Encode(Users[UsernameToId[unfollowedUsername]])
}
