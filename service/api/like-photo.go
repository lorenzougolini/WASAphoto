package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	photoid := uuid.Must(uuid.Parse(ps.ByName("photoid")))
	var likingPhoto Photo
	for _, photo := range Photos {
		if photo.photoID == photoid {
			likingPhoto = photo
			break
		}
	}

	// create new like object
	likeid, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// add check username exists

	newLike := Like{
		likeID:      likeid,
		username:    r.URL.Query().Get("username"),
		photoID:     photoid,
		dateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
	}
	likingPhoto.likes = append(likingPhoto.likes, newLike)
	json.NewEncoder(w).Encode(likingPhoto)
}
