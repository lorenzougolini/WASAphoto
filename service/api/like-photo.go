package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var message string
	userID := r.URL.Query().Get("userid")
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
	generateID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		newLikeID = generateID.String()
	}

	// create new Like object
	existsPhoto, likingPhoto, errP := rt.db.GetPhotoById(ps.ByName("photoid"))
	banned, errB := rt.db.IsBanned(likingPhoto.UserID, userID)
	if errP != nil || errB != nil {
		message = "Error liking the photo"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message)
		return

	} else if !existsPhoto {
		message = "Photo not found"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(message)
		return

	} else if banned {
		message = "User cannot like the photo"
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(message)
		return

	}

	newLike := Like{
		LikeID:      newLikeID,
		PhotoID:     likingPhoto.PhotoID,
		UserID:      userID,
		DateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
	}

	like, _ := json.Marshal(newLike)
	err = rt.db.AddLike(string(like))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(likingPhoto)
}
