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
	
	id, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id < 0 || id > len(Users){
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(Users[id])
	
}
