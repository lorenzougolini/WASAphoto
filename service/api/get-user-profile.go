package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"encoding/json"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if userID < 0 || userID > len(Users){
		w.WriteHeader(http.StatusBadRequest)
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if id not right token {
	// 	w.WriteHeader(http.StatuNotAuthorized)
	// 	return
	// }
	
	for i := range Users {
		if Users[i].username == username {
			user := Users[i]
			json.NewEncoder(w).Encode(user)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	}
