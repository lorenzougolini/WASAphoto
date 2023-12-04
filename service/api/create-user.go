package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) userLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	usernameString := r.URL.Query().Get("username")
	if usernameString == "" || len(usernameString) < 3 || len(usernameString) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return // check if uuid already exists
	}
	newUser := User{
		userID:   userID,
		username: usernameString,
		profile: Profile{
			photos:    []Photo{},
			followers: []User{},
			following: []User{},
			banned:    []User{},
		},
	}
	Users = append(Users, newUser)

	json.NewEncoder(w).Encode(newUser)
}
