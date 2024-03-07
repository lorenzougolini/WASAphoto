package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// check username to follow and proceed it exists
	followedUsername := ps.ByName("followedUsername")
	// fmt.Println(followedUsername)
	if followedUsername == username || followedUsername == "" || len(followedUsername) < 3 || len(followedUsername) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", followedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	followedUser, err := rt.db.GetByUsername(followedUsername)
	// fmt.Println(followedUser)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = fmt.Sprintf("The username '%s' doesn't exist", followedUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	err = rt.db.FollowUser(token, followedUser.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	message = followedUsername + " succesfully followed!"
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(message)
}
