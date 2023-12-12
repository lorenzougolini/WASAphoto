package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var message string
	userID := ps.ByName("userid")
	// check logged user id
	if !checkLogin(userID) {
		w.WriteHeader(http.StatusUnauthorized)
		message = "User is not correctly authenticated"
		json.NewEncoder(w).Encode(message)
		return
	} else if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" || len(username) < 3 || len(username) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", username)
		json.NewEncoder(w).Encode(message)
		return
	}

	// _, ok := UsernameToId[username]
	// if ok {
	// 	user, err := rt.db.GetUser(username)
	// 	if err != nil {
	// 		print(err)
	// 		return
	// 	} else {
	// 		message := user.UserID + ": " + user.Username
	// 		json.NewEncoder(w).Encode(message)
	// 	}
	// 	json.NewEncoder(w).Encode(Users[retrieveID])
	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	message = "Provided username does not exists"
	// 	json.NewEncoder(w).Encode(message)
	// 	return
	// }

	user, err := rt.db.GetUser(username)
	if err != nil {
		message = "Provided username does not exists"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message)
		return
	} else {
		// add call to getProfile
		message = user.UserID + ": " + user.Username
		json.NewEncoder(w).Encode(message)
	}

}
