package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) uploadNewPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	username := ps.ByName("username")
	if username != Logged.Username {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User is not authorized to add photos on this profile"
		json.NewEncoder(w).Encode(message)
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
	err = r.ParseForm()
	if err != nil {
		message = ("Failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	picture := r.Form.Get("picture")
	description := r.Form.Get("description")
	fmt.Printf("picture: %s, description: %s", picture, description)

	// create new Photo object
	newPhoto := Photo{
		PhotoID:     newPhotoID,
		UserID:      userID,
		Picture:     picture,
		DateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
		Description: description,
	}

	photo, _ := json.Marshal(newPhoto)
	err = rt.db.PostPhoto(string(photo))
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
