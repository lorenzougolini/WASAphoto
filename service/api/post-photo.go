package api

import (
	"encoding/json"
	"io"

	// "fmt"
	// "io"
	"net/http"
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

	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		message = ("Failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	picture, picHandler, err := r.FormFile("picture")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Error retrieving the file")
		return
	}
	defer picture.Close()

	// Check file size
	fileSize := picHandler.Size
	if fileSize > 3*1024*1024 { // Check if file size exceeds 3MB
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("File size exceeds size limit")
		return
	}

	// Opena and read image data
	file, err := picHandler.Open()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Error retrieving the file")
		return
	}
	defer file.Close()
	picData, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Error retrieving the file")
		return
	}

	description := r.FormValue("description")

	// create new Photo object
	newPhoto := Photo{
		PhotoID:     newPhotoID,
		AuthorID:    token,
		PicFile:     picData,
		DateAndTime: strconv.FormatInt(time.Now().Unix(), 10),
		Description: description,
	}

	photo, _ := json.Marshal(newPhoto)
	err = rt.db.PostPhoto(string(photo))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	_ = json.NewEncoder(w).Encode(newPhoto)
}
