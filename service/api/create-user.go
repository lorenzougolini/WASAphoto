package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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
	id, ok := UsernameToId[username]
	if ok {
		Logged["logged"] = Users[id].UserID
	} else {
		// user doesn't exists, so create a new one
		generateID, err := uuid.NewUUID()
		if _, ok := Users[generateID.String()]; err != nil || ok {
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
		UsernameToId[username] = newUserID
		Logged["logged"] = newUser.UserID
		Users[newUser.UserID] = newUser
		// u, _ := json.Marshal(newUser)
		// fmt.Println(string(u))
	}

	json.NewEncoder(w).Encode(Users[newUserID])
}
