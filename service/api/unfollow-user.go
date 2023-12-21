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
	userID := r.URL.Query().Get("userid")
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

	// check username to unfollow and proceed it exists
	unfollowedUsername := ps.ByName("unfollowedUsername")
	if (unfollowedUsername) == "" || len(string(unfollowedUsername)) < 3 || len(string(unfollowedUsername)) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", unfollowedUsername)
		json.NewEncoder(w).Encode(message)
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

	user, err := rt.db.GetByUsername(unfollowedUsername)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The username '%s' doesn't exist", unfollowedUsername)
		json.NewEncoder(w).Encode(message)
		return

	} else {
		err := rt.db.UnfollowUser(Logged.UserID, user.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		message = Logged.Username + " succesfully unfollowed: " + unfollowedUsername
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}

	json.NewEncoder(w).Encode(Users[UsernameToId[unfollowedUsername]])
}
