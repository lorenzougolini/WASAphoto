package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	var message string
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
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
		message = "No user found"
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	if banned, err := rt.db.IsBannedBy(token, user.UserID); banned || err != nil {

		message = "The user is banned"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	profile, err := rt.db.GetProfile(user.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	profile.Banned, err = rt.db.IsBannedBy(user.UserID, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	profileJson, err := json.MarshalIndent(profile, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	_, _ = w.Write(profileJson)
}
