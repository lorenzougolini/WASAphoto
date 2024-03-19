package api

import (
	"encoding/json"
	"net/http"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	var message string
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
		return
	}

	photoID := ps.ByName("photoid")
	likeID := ps.ByName("likeid")
	existsPhoto, _, errP := rt.db.GetPhotoById(photoID)
	existsLike, removedLike, errL := rt.db.GetLikeByLikeId(likeID)
	if !existsPhoto {
		w.WriteHeader(http.StatusNotFound)
		message = "Photo not found"
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	if !existsLike {
		w.WriteHeader(http.StatusNotFound)
		message = "Like not found"
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	if errP != nil || errL != nil {
		message = "Error removing like"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	err := rt.db.RemoveLike(removedLike.LikeID)
	if err != nil {
		message = "Error removing like"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	unlikedPhotoData, _ := rt.db.GetPhotoData(photoID)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(unlikedPhotoData)
}
