package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) uploadNewPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "multipart/form-data")

	var message string
	username := ps.ByName("username")
	// check Bearer token
	if !checkLogin(r) || username != Logged.Username {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}
	if username != Logged.Username {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User is not authorized to add photos on this profile"
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	var newPhotoID string
	generateID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		newPhotoID = generateID.String()
	}

	// SECOND METHOD TO READ BODY
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		message = ("Failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	picture, picHandler, err := r.FormFile("picture")
	if err != nil {
		_ = json.NewEncoder(w).Encode("Error retrieving the file")
		return
	}
	defer picture.Close()
	description := r.FormValue("description")
	// fmt.Printf("picture: %s, description: %s", picture, description)

	picDir := "./pictures/"
	if _, err := os.Stat(picDir); os.IsNotExist(err) {
		err := os.Mkdir(picDir, os.ModePerm)
		if err != nil {
			_ = json.NewEncoder(w).Encode("Error creating upload directory")
			return
		}
	}

	// create new Photo object
	newPhoto := Photo{
		PhotoID:     newPhotoID,
		UserID:      Logged.UserID,
		PicPath:     picDir + newPhotoID + ".jpg",
		DateAndTime: time.Now().Format("2017-07-21T17:32:28Z"),
		Description: description,
	}

	dst, err := os.Create(newPhoto.PicPath)
	if err != nil {
		_ = json.NewEncoder(w).Encode("Error creating the file")
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, picture); err != nil {
		_ = json.NewEncoder(w).Encode("Error copying the file")
		return
	}
	// Respond with a success message
	message = fmt.Sprintf("File %s uploaded successfully!", picHandler.Filename)
	_ = json.NewEncoder(w).Encode(message)

	photo, _ := json.Marshal(newPhoto)
	err = rt.db.PostPhoto(string(photo))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	// Photos[newPhotoID] = newPhoto
	// getUser := Users[userID]
	// getUser.Profile.photos = append(getUser.Profile.photos, newPhoto.PhotoID)
	_ = json.NewEncoder(w).Encode(newPhoto)
}
