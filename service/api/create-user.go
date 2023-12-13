package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "github.com/google/uuid"
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
		json.NewEncoder(w).Encode(message)
		return
	}

	// check if username already exists, if yes log him in
	var newUserID string
	user, err := rt.db.GetUser(username)
	if err != nil {
		print(err)
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
			Profile: Profile{
				photos:    []string{},
				followers: []string{},
				following: []string{},
				banned:    []string{},
			},
		}
		// UsernameToId[username] = newUserID

		Logged["id"] = newUser.UserID
		Logged["username"] = newUser.Username

		// Users[newUser.UserID] = newUser
		err = rt.db.SetUser(newUserID, newUser.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
	} else {
		Logged["id"] = user.UserID
		Logged["username"] = user.Username
	}

	json.NewEncoder(w).Encode(Logged)
}
