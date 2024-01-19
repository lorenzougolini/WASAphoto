package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var message string
	username := ps.ByName("username")
	// check Bearer token
	if !checkLogin(r) || username != Logged.Username {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	// check username to follow and proceed it exists
	followedUsername := ps.ByName("followedUsername")
	// fmt.Println(followedUsername)
	if followedUsername == username || followedUsername == "" || len(followedUsername) < 3 || len(followedUsername) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", followedUsername)
		_ = json.NewEncoder(w).Encode(message)
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

	followedUser, err := rt.db.GetByUsername(followedUsername)
	// fmt.Println(followedUser)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The username '%s' doesn't exist", followedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	} else {
		logged, _ := json.Marshal(Logged)
		followed, _ := json.Marshal(followedUser)
		err := rt.db.FollowUser(string(logged), string(followed))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
		message = Logged.Username + " succesfully followed: " + followedUsername
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(message)
	}

	_ = json.NewEncoder(w).Encode(followedUser)
}
