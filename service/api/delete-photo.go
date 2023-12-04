package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")

	var message string
	pathUserID := ps.ByName("userid")
	if pathUserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User is not correctely authenticated"
		json.NewEncoder(w).Encode(message)
	}

	// find the user and the photo to remove
	pathPhotoID := ps.ByName("photoid")
	var loggedUser User
	for _, user := range Users {
		if user.userID == userID {
			loggedUser = user
			break
		}
	}

	var ind int
	for i, photo := range loggedUser.profile.photos {
		if photo.photoID == uuid.Must(uuid.Parse(pathPhotoID)) {
			ind = i
			break
		}
	}

	loggedUser.profile.photos = append(loggedUser.profile.photos[:ind], loggedUser.profile.photos[ind+1:]...)
	json.NewEncoder(w).Encode(loggedUser)
}
