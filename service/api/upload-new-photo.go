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
func (rt *_router) uploadNewPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")

	var message string
	pathUserID := ps.ByName("userid")
	if pathUserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User is not correctely authenticated"
		json.NewEncoder(w).Encode(message)
	}
	var loggedUser User
	for _, user := range Users {
		if user.userID == userID {
			loggedUser = user
			break
		}
	}

	// assign photoID
	photoid, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// retrieve picture from request body
	picture, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		message = "Server is unable to process the uploaded file"
		json.NewEncoder(w).Encode(message)
	}

	// create new Photo object
	newPhoto := Photo{
		photoID:     photoid,
		userID:      uuid.Must(uuid.Parse(pathUserID)),
		picture:     string(picture),
		dateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
		likes:       []Like{},
		comments:    []Comment{},
	}
	loggedUser.profile.photos = append(loggedUser.profile.photos, newPhoto)
	json.NewEncoder(w).Encode(newPhoto)

}
