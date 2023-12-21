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
		message = "Server is unable to process the request body"
		json.NewEncoder(w).Encode(message)
		return
	}

	commentText := string(text)
	if len(text) < 1 || len(text) > 100 {
		w.WriteHeader(http.StatusBadRequest)
		message = "The submitted comment is not valid"
		json.NewEncoder(w).Encode(message)
		return
	}

	// create new comment object
	existsPhoto, commentingPhoto, err := rt.db.GetPhotoById(ps.ByName("photoid"))
	if !existsPhoto {
		message = "Photo not found"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(message)
		return
	}

	newComment := Comment{
		CommentID:   newCommentID,
		UserID:      userID,
		PhotoID:     commentingPhoto.PhotoID,
		CommentText: commentText,
		DateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
	}

	comment, _ := json.Marshal(newComment)
	err = rt.db.AddComment(string(comment))
	json.NewEncoder(w).Encode(commentingPhoto)

}
