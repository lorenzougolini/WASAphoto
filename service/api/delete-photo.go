package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// delete the photo from Photos map and User profile

	photoID := ps.ByName("photoid")
	if _, ok := Photos[photoID]; ok {
		getUser := Users[userID]
		getUser.Profile.photos = remove(getUser.Profile.photos, photoID)
		delete(Photos, photoID)
	} else {
		w.WriteHeader(http.StatusNotFound)
		message = "Server unable to find the photo"
		json.NewEncoder(w).Encode(message)
	}

	json.NewEncoder(w).Encode(Logged["logged"])
}