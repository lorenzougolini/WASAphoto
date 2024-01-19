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
	// check Bearear token
	if !checkLogin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	// get comment text from request body
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		message = ("Failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	commentText := r.FormValue("comment")

	if len(commentText) < 1 || len(commentText) > 100 {
		w.WriteHeader(http.StatusBadRequest)
		message = "The submitted comment is not valid"
		_ = json.NewEncoder(w).Encode(message)
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

	existsPhoto, commentingPhoto, errP := rt.db.GetPhotoById(ps.ByName("photoid"))
	banned, errB := rt.db.IsBanned(commentingPhoto.UserID, Logged.UserID)
	if errP != nil || errB != nil {
		message = "Error commenting the photo"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return

	} else if !existsPhoto {
		message = "Photo not found"
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(message)
		return

	} else if banned {
		message = "User cannot like the photo"
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	newComment := Comment{
		CommentID:   newCommentID,
		UserID:      Logged.UserID,
		PhotoID:     commentingPhoto.PhotoID,
		CommentText: commentText,
		DateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
	}

	comment, _ := json.Marshal(newComment)
	err = rt.db.AddComment(string(comment))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	_ = json.NewEncoder(w).Encode(commentingPhoto)

}
