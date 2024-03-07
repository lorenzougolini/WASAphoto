package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	var message string
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
		return
	}

	// check username to follow and proceed it exists
	reqUsername := ps.ByName("username")
	unbannedUsername := ps.ByName("unbannedUsername")
	if unbannedUsername == reqUsername || unbannedUsername == "" || len(unbannedUsername) < 3 || len(unbannedUsername) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", unbannedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	unbannedUser, err := rt.db.GetByUsername(unbannedUsername)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The username '%s' doesn't exist", unbannedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return

	}

	err = rt.db.UnbanUser(token, unbannedUser.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	message = unbannedUser.Username + " succesfully unbanned!"
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(message)

}
