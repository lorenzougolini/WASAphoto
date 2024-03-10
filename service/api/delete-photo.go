package api

import (
	"encoding/json"
	"net/http"
	"os"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	var message string
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
		return
	}

	// delete the photo from Photos map and User profile
	photoID := ps.ByName("photoid")
	photoExists, retrievedPhoto, err := rt.db.GetPhotoById(photoID)
	if !photoExists || err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = "Server unable to find the photo"
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	if retrievedPhoto.AuthorID != token {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User in not authorized to remove photos from this profile"
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	err = rt.db.RemovePhoto(retrievedPhoto.PhotoID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	photoDir := "./service/imgDB/" + retrievedPhoto.PhotoID + ".jpg"
	err = os.Remove(photoDir)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	message = "Photo removed"
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(message)
}
