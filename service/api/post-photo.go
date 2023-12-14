package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
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
	generateID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		newPhotoID = generateID.String()
	}

	// FIRST METHOD TO READ BODY
	// picture, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnprocessableEntity)
	// 	message = "Server is unable to process the uploaded file"
	// 	json.NewEncoder(w).Encode(message)
	// 	return
	// }

	// SECOND METHOD TO READ BODY
	type Picture struct {
		Picture     string // int64
		Description string
	}
	var picture Picture
	err = json.NewDecoder(r.Body).Decode(&picture)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message = "Impossible to read the request body"
		json.NewEncoder(w).Encode(message)
		return
	}

	// create new Photo object
	newPhoto := Photo{
		PhotoID:     newPhotoID,
		UserID:      userID,
		Picture:     string(picture.Picture),
		DateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
		Description: string(picture.Description),
	}

	err = rt.db.PostPhoto(newPhoto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	// Photos[newPhotoID] = newPhoto
	// getUser := Users[userID]
	// getUser.Profile.photos = append(getUser.Profile.photos, newPhoto.PhotoID)
	json.NewEncoder(w).Encode(newPhoto)
}
