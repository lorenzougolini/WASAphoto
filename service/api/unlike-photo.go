package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var message string
	userID := ps.ByName("userid")
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

	photoID := ps.ByName("photoid")
	likeID := ps.ByName("likeid")
	_, okp := Photos[photoID]
	_, okl := Likes[likeID]
	if !okp || !okl {
		w.WriteHeader(http.StatusNotFound)
		message = "Like or photo not found"
		json.NewEncoder(w).Encode(message)
		return
	} else {
		// remove like from list and from photo
	}

	likingPhoto.likes = append(likingPhoto.likes[:ind], likingPhoto.likes[ind+1:]...)
	json.NewEncoder(w).Encode(likingPhoto)
}
