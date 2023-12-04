package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	// "strconv"
)

var Users = []User{}
var Photos = []Photo{}

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) listAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	if item := r.URL.Query().Get("item"); item == "users" {
		json.NewEncoder(w).Encode(Users)
	} else {
		json.NewEncoder(w).Encode(Photos)
	}
}
