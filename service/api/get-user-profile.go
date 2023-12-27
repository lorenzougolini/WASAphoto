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
	userID := r.URL.Query().Get("userid")
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

	username := ps.ByName("username")
	if username == "" || len(username) < 3 || len(username) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", username)
		json.NewEncoder(w).Encode(message)
		return
	}

	user, err := rt.db.GetByUsername(username)
	if err != nil {
		message = "Provided username does not exists"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message)
		return

	} else {
		// add call to getProfile
		profile, err := rt.db.GetProfile(user.UserID)
		if err != nil {
			message = "Error retrieving the profile"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(message)
			return
		}

		profileJson, err := json.MarshalIndent(profile, "", " ")
		if err != nil {
			message = "Error retrieving the profile"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(message)
			return
		}

		message = user.Username
		json.NewEncoder(w).Encode(message)
		fmt.Println(string(profileJson))
	}

}
