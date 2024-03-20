package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "multipart/form-data")

	var message string
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
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
		newCommentID = formatId(generateID.String())
	}

	existsPhoto, commentingPhoto, errP := rt.db.GetPhotoById(ps.ByName("photoid"))
	if errP != nil {
		message = "Error commenting the photo"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return

	} else if !existsPhoto {
		message = "Photo not found"
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	banned, errB := rt.db.IsBannedBy(token, commentingPhoto.AuthorID)
	if banned || errB != nil {
		message = "User cannot like the photo"
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	newComment := Comment{
		CommentID:   newCommentID,
		PhotoID:     commentingPhoto.PhotoID,
		User:        token,
		CommentText: commentText,
		DateAndTime: strconv.FormatInt(time.Now().Unix(), 10),
	}

	comment, _ := json.Marshal(newComment)
	err = rt.db.AddComment(string(comment))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	// ctx.Logger.Info("Comment added to photo " + newComment.PhotoID + ": " + newComment.CommentText)
	commentedPhoto, _ := rt.db.GetPhotoData(newComment.PhotoID)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(commentedPhoto)

}
