package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) userLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// check if username already exists, if yes log him in
	var newUserID string
	user, err := rt.db.GetByUsername(username)
	if err != nil {
		// user doesn't exists, so create a new one
		generateID, err := uuid.NewV4()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if exists, err := rt.db.CheckID(generateID.String()); err != nil || exists > 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			newUserID = generateID.String()
		}
		newUser := User{
			UserID:   newUserID,
			Username: username,
		}
		// UsernameToId[username] = newUserID

		Logged.UserID = newUser.UserID
		Logged.Username = newUser.Username
		// rt.usr.Username = newUser.Username

		// Users[newUser.UserID] = newUser
		err = rt.db.SetUser(newUserID, newUser.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
	} else {
		Logged.UserID = user.UserID
		Logged.Username = user.Username
	}

	_ = json.NewEncoder(w).Encode(Logged)
}
