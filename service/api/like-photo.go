package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// create new LikeID
	var newLikeID string
	generateID, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		newLikeID = generateID.String()
	}

	// create new Like object
	likingPhoto := Photos[ps.ByName("photoid")]
	newLike := Like{
		likeID:      newLikeID,
		username:    Users[userID].Username,
		photoID:     likingPhoto.PhotoID,
		dateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
	}

	Likes[newLikeID] = newLike
	likingPhoto.Likes = append(likingPhoto.Likes, newLike.username)
	json.NewEncoder(w).Encode(likingPhoto)
}
