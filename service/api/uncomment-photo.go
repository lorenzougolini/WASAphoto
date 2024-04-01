package api

import (
	"encoding/json"
	"net/http"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "text/plain")

	var message string
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
		return
	}

	photoID := ps.ByName("photoid")
	commentID := ps.ByName("commentid")
	existsPhoto, _, errP := rt.db.GetPhotoById(photoID)
	existsComment, removedComment, errC := rt.db.GetCommentByCommentId(commentID)
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
		message = "Error removing comment"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return

	}

	err := rt.db.RemoveComment(removedComment.CommentID)
	if err != nil {
		message = "Error removing comment"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	uncommentedPhotoData, _ := rt.db.GetPhotoData(photoID)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(uncommentedPhotoData)
}
