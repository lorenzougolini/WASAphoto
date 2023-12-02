package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type Photo struct {
	photoID     uuid.UUID
	userID      uuid.UUID
	picture     string
	dateAndTime time.Time
	likes       []Like
	comments    []Comment
}

type Like struct {
	likeID      uuid.UUID
	username    string
	photoID     uuid.UUID
	dateAndTime time.Time
}

type Comment struct {
	commentID   uuid.UUID
	username    string
	photoID     uuid.UUID
	commentText string
	dateAndTime time.Time
}

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) uploadNewPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	id := 2
	json.NewEncoder(w).Encode(id)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// if userID < 0 || userID > len(Users){
	//  	w.WriteHeader(http.StatusBadRequest)
	// }
	// add check if userID incorrect

	// if id not right token {
	// 	w.WriteHeader(http.StatuNotAuthorized)
	// 	return
	// }
	// id, err := uuid.NewUUID()
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// }
	// newPhoto := Photo{
	// 	photoID: id,
	// 	userID:  userID,
	// 	picture: "0",
	// }

}
