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
}
