package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "multipart/form-data")

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

	// get comment text from request body
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		message = ("Failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	commentText := r.FormValue("comment")

	if len(commentText) < 1 || len(commentText) > 100 {
		w.WriteHeader(http.StatusBadRequest)
		message = "The submitted comment is not valid"
		json.NewEncoder(w).Encode(message)
		return
	}

	// create new comment object
	var newCommentID string
	generateID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		newCommentID = generateID.String()
	}

	existsPhoto, commentingPhoto, err := rt.db.GetPhotoById(ps.ByName("photoid"))
	if err != nil {
		message = "Error commenting the photo"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message)
		return

	} else if !existsPhoto {
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(commentingPhoto)

}
