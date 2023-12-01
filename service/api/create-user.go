package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	// "github.com/google/uuid"
)

type User struct {
	userID int
	username string
	profile Profile //change to Profile
}

type Profile struct {
	photos []int // change to photo
	followers []User
	following []User
	banned []User
}

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := len(Users)
	usernameString := r.URL.Query().Get("username")
	if usernameString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser := User{
		userID: id,
		username: usernameString,
		profile: 0,
	}
	Users = append(Users, newUser)
	response := map[string]string{"message": "User created successfully", "username": usernameString}

	json.NewEncoder(w).Encode(response)
}
