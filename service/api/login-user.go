package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"WASAphoto.uniroma1.it/WASAphoto/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) userLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	username := strings.TrimSpace(r.URL.Query().Get("username"))
	var message string
	// check valid username
	if username == "" || len(username) < 3 || len(username) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided username '%s' is not valid", username)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	// check if username already exists, if yes log in
	var newUserID string

	user, err := rt.db.GetByUsername(username)

	if err != nil {

		// user doesn't exists, create a new one
		generateID, _ := uuid.NewV4()
		newUserID = formatId(generateID.String())
		if exists, err := rt.db.CheckIDExistence(newUserID); err != nil || exists {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = rt.db.SetUser(newUserID, username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
		ctx.Logger.Info("User created")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(User{UserID: newUserID, Username: username})

	} else {

		ctx.Logger.Info("User logged in")
		_ = json.NewEncoder(w).Encode(user)
	}

}
