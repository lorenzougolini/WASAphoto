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
	username := ps.ByName("username")
	// check Bearer token
	if !checkLogin(r) || username != Logged.Username {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	// check username to follow and proceed it exists
	bannedUsername := ps.ByName("bannedUsername")
	if bannedUsername == username || bannedUsername == "" || len(bannedUsername) < 3 || len(bannedUsername) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", bannedUsername)
		_ = json.NewEncoder(w).Encode(message)
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
		_ = json.NewEncoder(w).Encode(message)
		return
	} else {
		logged, _ := json.Marshal(Logged)
		err1 := rt.db.BanUser(string(logged), user.Username)
		err2 := rt.db.UnfollowUser(Logged.UserID, user.Username)
		if err1 != nil || err2 != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
		message = Logged.Username + " succesfully banned: " + user.Username
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(message)
	}

	// _ = json.NewEncoder(w).Encode(Users[UsernameToId[bannedUsername]])
}
