package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	photoid := uuid.Must(uuid.Parse(ps.ByName("photoid")))
	likeid := uuid.Must(uuid.Parse(ps.ByName("likeid")))

	var likingPhoto Photo
	for _, photo := range Photos {
		if photo.photoID == photoid {
			likingPhoto = photo
			break
		}
	}

	// find like object
	var ind int
	for i, like := range likingPhoto.likes {
		if like.likeID == likeid {
			ind = i
			break
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	likingPhoto.likes = append(likingPhoto.likes[:ind], likingPhoto.likes[ind+1:]...)
	json.NewEncoder(w).Encode(likingPhoto)
}
