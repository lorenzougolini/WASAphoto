package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	var message string
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
		return
	}

	// create new LikeID
	var newLikeID string
	generateID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		newLikeID = formatId(generateID.String())
	}

	// create new Like object
	existsPhoto, likingPhoto, errP := rt.db.GetPhotoById(ps.ByName("photoid"))
	banned, errB := rt.db.IsBannedBy(token, likingPhoto.AuthorID)
	if errP != nil || errB != nil {
		message = "Error liking the photo"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	if !existsPhoto {
		message = "Photo not found"
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	if banned {
		message = "User cannot like the photo"
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	// check if the user has already liked the photo
	existLike, errL := rt.db.GetLikeByUserId(token, likingPhoto.PhotoID)
	if errL != nil {
		message = "Error liking the photo"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	if existLike {
		message = "User has already liked the photo"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	newLike := Like{
		LikeID:      newLikeID,
		PhotoID:     likingPhoto.PhotoID,
		UserID:      token,
		DateAndTime: strconv.FormatInt(time.Now().Unix(), 10),
	}

	like, _ := json.Marshal(newLike)
	err = rt.db.AddLike(string(like))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	likedPhotoData, _ := rt.db.GetPhotoData(ps.ByName("photoid"))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(likedPhotoData)
}
