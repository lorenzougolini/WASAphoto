package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var message string
	userID := r.URL.Query().Get("userid")
	username := ps.ByName("username")
	// check logged user id
	if !checkLogin(userID) || username != Logged.Username {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User is not correctly authenticated"
		json.NewEncoder(w).Encode(message)
		return
	} else if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check username to follow and proceed it exists
	bannedUsername := ps.ByName("bannedUsername")
	if bannedUsername == username || bannedUsername == "" || len(string(bannedUsername)) < 3 || len(string(bannedUsername)) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", bannedUsername)
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

	user, err := rt.db.GetByUsername(bannedUsername)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The username '%s' doesn't exist", bannedUsername)
		json.NewEncoder(w).Encode(message)
		return
	} else {
		logged, _ := json.Marshal(Logged)
		err1 := rt.db.BanUser(string(logged), user.Username)
		err2 := rt.db.UnfollowUser(Logged.UserID, user.Username)
		if err1 != nil || err2 != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		message = Logged.Username + " succesfully banned: " + user.Username
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}

	json.NewEncoder(w).Encode(Users[UsernameToId[bannedUsername]])
}
