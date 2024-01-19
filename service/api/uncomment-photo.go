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
	// check Bearer token
	if !checkLogin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	photoID := ps.ByName("photoid")
	commentID := ps.ByName("commentid")
	existsPhoto, uncommentedPhoto, errP := rt.db.GetPhotoById(photoID)
	existsComment, removedComment, errC := rt.db.GetCommentById(commentID)
	if !existsPhoto {
		w.WriteHeader(http.StatusNotFound)
		message = "Photo not found"
		_ = json.NewEncoder(w).Encode(message)
		return

	} else if !existsComment {
		w.WriteHeader(http.StatusNotFound)
		message = "Comment not found"
		_ = json.NewEncoder(w).Encode(message)
		return

	} else if errP != nil || errC != nil {
		message = "Error removing like"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return

	} else {
		err := rt.db.RemoveComment(removedComment.CommentID, uncommentedPhoto.PhotoID)
		if err != nil {
			message = "Error removing comment"
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(message)
			return
		}
	}

	_ = json.NewEncoder(w).Encode(uncommentedPhoto)
}
