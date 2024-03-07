package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) userLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	username := r.URL.Query().Get("username")
	var message string
	// check valid username
	if username == "" || len(username) < 3 || len(username) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", username)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	// check if username already exists, if yes log in
	var newUserID string

	dbuser, err := rt.db.GetByUsername(username)
	if err != nil {

		// user doesn't exists, create a new one
		generateID, err := uuid.NewV4()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newUserID = generateID.String()
		if exists, err := rt.db.CheckIDExistence(newUserID); err != nil || exists {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newUser := User{
			UserID:   newUserID,
			Username: username,
		}

		dbuser, err = rt.db.SetUser(newUser.userToDBUser())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err)
			return
		}

	}

	_ = json.NewEncoder(w).Encode(dbuser)
}
