package api

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/oknors/okno/app/mod"
	"net/http"
)

// Create appends post path prefix for a database write
func (a *API) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["host"]
	col := mux.Vars(r)["col"]
	id := mux.Vars(r)["slug"]
	err := r.ParseForm()
	if err != nil {
		// Handle error
	}
	var post mod.Post
	// r.PostForm is a map of our POST form values
	err = decoder.Decode(&post, r.PostForm)
	if err != nil {
		// Handle error
	}
	a.Write(path+"/"+col, id, post)

}

var decoder = schema.NewDecoder()

// Change host of API
func (a *API) Host(h string) *API {
	a.path = h
	return a
}
