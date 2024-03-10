package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "text/plain")

	var message string
	// check Bearer token
	token := r.Header.Get("Authorization")
	if exists, err := rt.db.CheckIDExistence(token); err != nil || token == "" || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errUncorrectLogin)
		return
	}

	newUsername, err := io.ReadAll(r.Body)
	if string(newUsername) == "" || len(string(newUsername)) < 3 || len(string(newUsername)) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", newUsername)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		message = "Error changing the username"
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	err = rt.db.SetName(token, string(newUsername))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info("The username has been changed")
	_ = json.NewEncoder(w).Encode(User{UserID: token, Username: string(newUsername)})
}
