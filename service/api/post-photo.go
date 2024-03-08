package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) uploadNewPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "multipart/form-data")

	var message string
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
		return
	}

	var newPhotoID string
	generateID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		newPhotoID = formatId(generateID.String())
	}

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

	picDir := "./webui/public/pictures/"
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
		AuthorID:    token,
		PicPath:     picDir + newPhotoID + ".jpg",
		DateAndTime: strconv.FormatInt(time.Now().Unix(), 10),
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

	_ = json.NewEncoder(w).Encode(newPhoto)
}
