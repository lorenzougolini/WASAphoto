package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type User struct{
	userID uuid
	username string
	profile Profile
}

type Profile struct{
	photos []Photo
	followers []User
	follwing []User
	bannedUsers []User
}

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Hello World!\n"))
}
