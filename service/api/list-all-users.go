package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	// "strconv"
)

var Users = []User{}

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) listAllUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	// id, err := strconv.Atoi(ps.ByName("id"))
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	json.NewEncoder(w).Encode(Users)
}
