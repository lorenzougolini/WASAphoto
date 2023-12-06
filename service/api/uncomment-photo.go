package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	photoID := ps.ByName("photoid")
	commentID := ps.ByName("commentid")
	uncommentingPhoto, okP := Photos[photoID]
	_, okC := Comments[commentID]

	if !okP || !okC {
		w.WriteHeader(http.StatusNotFound)
		message = "Comment or photo not found"
		json.NewEncoder(w).Encode(message)
	} else {
		// remove comment from map and from associated photo
		uncommentingPhoto.Comments = remove(uncommentingPhoto.Comments, commentID)
		delete(Comments, commentID)
	}

	json.NewEncoder(w).Encode(uncommentingPhoto)
}
