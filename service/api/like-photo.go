package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var message string
	// check Bearer token
	if !checkLogin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	// create new LikeID
	var newLikeID string
	generateID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		newLikeID = generateID.String()
	}

	// create new Like object
	existsPhoto, likingPhoto, errP := rt.db.GetPhotoById(ps.ByName("photoid"))
	banned, errB := rt.db.IsBanned(likingPhoto.UserID, Logged.UserID)
	if errP != nil || errB != nil {
		message = "Error liking the photo"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return

	} else if !existsPhoto {
		message = "Photo not found"
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(message)
		return

	} else if banned {
		message = "User cannot like the photo"
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(message)
		return

	}

	// check if the user has already liked the photo
	existLike, errL := rt.db.GetLikeByUserId(Logged.UserID, likingPhoto.PhotoID)
	if errL != nil {
		message = "Error liking the photo"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return
	} else if existLike {
		message = "User has already liked the photo"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	newLike := Like{
		LikeID:      newLikeID,
		PhotoID:     likingPhoto.PhotoID,
		UserID:      Logged.UserID,
		DateAndTime: strconv.FormatInt(time.Now().Unix(), 10),
	}

	like, _ := json.Marshal(newLike)
	err = rt.db.AddLike(string(like))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	_ = json.NewEncoder(w).Encode(likingPhoto)
}
