package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var message string
	// check Bearer token
	if !checkLogin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	photoID := ps.ByName("photoid")
	likeID := ps.ByName("likeid")
	existsPhoto, unlikedPhoto, errP := rt.db.GetPhotoById(photoID)
	existsLike, removedLike, errL := rt.db.GetLikeById(likeID)
	if !existsPhoto {
		w.WriteHeader(http.StatusNotFound)
		message = "Photo not found"
		_ = json.NewEncoder(w).Encode(message)
		return

	} else if !existsLike {
		w.WriteHeader(http.StatusNotFound)
		message = "Like not found"
		_ = json.NewEncoder(w).Encode(message)
		return

	} else if errP != nil || errL != nil {
		message = "Error removing like"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return

	} else {
		err := rt.db.RemoveLike(removedLike.LikeID, unlikedPhoto.PhotoID)
		if err != nil {
			message = "Error removing like"
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(message)
			return
		}
	}

	_ = json.NewEncoder(w).Encode(unlikedPhoto)
}
