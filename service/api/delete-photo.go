package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var message string
	username := ps.ByName("username")
	// check Bearer token
	if !checkLogin(r) || username != Logged.Username {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	// delete the photo from Photos map and User profile
	photoID := ps.ByName("photoid")
	photoExists, retrievedPhoto, err := rt.db.GetPhotoById(photoID)
	if !photoExists || err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = "Server unable to find the photo"
		_ = json.NewEncoder(w).Encode(message)

	} else {
		if retrievedPhoto.UserID != Logged.UserID {
			w.WriteHeader(http.StatusUnauthorized)
			message = "User in not authorized to remove photos from this profile"
			_ = json.NewEncoder(w).Encode(message)
			return

		} else {
			err = rt.db.RemovePhoto(Logged.UserID, retrievedPhoto.PhotoID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(err)
				return
			}
		}
	}

	// if _, ok := Photos[photoID]; ok {
	// 	getUser := Users[userID]
	// 	getUser.Profile.photos = remove(getUser.Profile.photos, photoID)
	// 	delete(Photos, photoID)
	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	message = "Server unable to find the photo"
	// 	json.NewEncoder(w).Encode(message)
	// }

	_ = json.NewEncoder(w).Encode(Logged)
}
