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
	unlikedPhoto, okP := Photos[photoID]
	_, okL := Likes[likeID]
	if !okP || !okL {
		w.WriteHeader(http.StatusNotFound)
		message = "Like or photo not found"
		json.NewEncoder(w).Encode(message)
		return
	} else {
		// remove like from map and from associated photo
		unlikedPhoto.Likes = remove(unlikedPhoto.Likes, Users[userID].Username)
		delete(Likes, likeID)
	}

	json.NewEncoder(w).Encode(unlikedPhoto)
}
