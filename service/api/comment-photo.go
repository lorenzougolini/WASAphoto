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
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var newCommentID string
	generateID, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		newCommentID = generateID.String()
	}

	// get comment text from request body
	text, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		message = "Server is unable to process the uploaded file"
		json.NewEncoder(w).Encode(message)
		return
	}
	commentText := string(text)
	if len(text) < 1 || len(text) > 100 {
		w.WriteHeader(http.StatusBadRequest)
		message = "The comment is out of bounds"
		json.NewEncoder(w).Encode(message)
		return
	}

	// create new comment object
	commentingPhoto := Photos[ps.ByName("photoid")]
	newComment := Comment{
		commentID:   newCommentID,
		username:    Users[userID].Username,
		photoID:     commentingPhoto.PhotoID,
		commentText: commentText,
		dateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
	}

	Comments[newCommentID] = newComment
	commentingPhoto.Comments = append(commentingPhoto.Comments, newCommentID)
	json.NewEncoder(w).Encode(commentingPhoto)

}
