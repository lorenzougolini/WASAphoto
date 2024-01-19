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
	// check Bearer token
	if !checkLogin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(uncorrectLogin)
		return
	}

	username := ps.ByName("username")
	if username == "" || len(username) < 3 || len(username) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", username)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	user, err := rt.db.GetByUsername(username)
	if err != nil {
		message = "Provided username does not exists"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return

	} else {
		// add call to getProfile
		profile, err := rt.db.GetProfile(user.UserID)
		if err != nil {
			message = "Error retrieving the profile"
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(message)
			return
		}

		profileJson, err := json.MarshalIndent(profile, "", " ")
		if err != nil {
			message = "Error retrieving the profile"
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(message)
			return
		}

		message = user.Username
		_ = json.NewEncoder(w).Encode(message)
		w.Write(profileJson)
		// logrus.Println(string(profileJson))
	}

}
