package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	var message string
	username := ps.ByName("username")
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
		return
	}

	// check username to ban and proceed it exists
	bannedUsername := ps.ByName("bannedUsername")
	if bannedUsername == username || bannedUsername == "" || len(bannedUsername) < 3 || len(bannedUsername) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", bannedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	bannedUser, err := rt.db.GetByUsername(bannedUsername)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The username '%s' doesn't exist", bannedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	err = rt.db.BanUser(token, bannedUser.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	err = rt.db.UnfollowUser(token, bannedUser.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	err = rt.db.UnfollowUser(bannedUser.UserID, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	message = bannedUser.Username + " succesfully banned!"
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(message)

}
