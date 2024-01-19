package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	unbannedUsername := ps.ByName("unbannedUsername")
	if unbannedUsername == username || unbannedUsername == "" || len(unbannedUsername) < 3 || len(unbannedUsername) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", unbannedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	// if _, ok := UsernameToId[username]; ok {

	// 	// add folllowing
	// 	getUser := Users[userID]

	// 	if ok2 := contains(getUser.Profile.banned, username); ok2 {
	// 		getUser.Profile.banned = remove(getUser.Profile.banned, username)
	// 	} else {
	// 		w.WriteHeader(http.StatusNotFound)
	// 		message = fmt.Sprintf("The user '%s' is not banned", username)
	// 		json.NewEncoder(w).Encode(message)
	// 		return
	// 	}

	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	message = fmt.Sprintf("The user '%s' doesn't exist", username)
	// 	json.NewEncoder(w).Encode(message)
	// 	return
	// }

	user, err := rt.db.GetByUsername(unbannedUsername)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The username '%s' doesn't exist", unbannedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return

	} else {
		err := rt.db.UnbanUser(Logged.UserID, user.UserID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
		message = Logged.Username + " succesfully unbanned: " + user.Username
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(message)
	}

	// _ = json.NewEncoder(w).Encode(Users[UsernameToId[unbannedUsername]])
}
