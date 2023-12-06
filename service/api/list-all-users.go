package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	// "strconv"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) listAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	for _, user := range Users {
		json.NewEncoder(w).Encode(user)
	}
}
