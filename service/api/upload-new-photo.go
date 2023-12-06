package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) uploadNewPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")

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

	var newPhotoID string
	generateID, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		newPhotoID = generateID.String()
	}

	// retrieve picture from request body
	picture, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		message = "Server is unable to process the uploaded file"
		json.NewEncoder(w).Encode(message)
		return
	}

	// create new Photo object
	newPhoto := Photo{
		PhotoID:     newPhotoID,
		UserID:      userID,
		Picture:     string(picture),
		DateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
		Likes:       []string{},
		Comments:    []string{},
	}

	Photos[newPhotoID] = newPhoto
	getUser := Users[userID]
	getUser.Profile.photos = append(getUser.Profile.photos, newPhoto.PhotoID)
	json.NewEncoder(w).Encode(newPhoto)
}
