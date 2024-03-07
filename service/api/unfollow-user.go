package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "text/plain")

	var message string
	username := ps.ByName("username")
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
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
		message = fmt.Sprintf("The user '%s' doesn't exist", unfollowedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	err = rt.db.UnfollowUser(token, unfollowedUser.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	message = unfollowedUsername + " succesfully unfollowed!"
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(message)

}
