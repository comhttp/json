package jdb

import (
	"github.com/gorilla/schema"

	"github.com/gorilla/mux"
	"github.com/oknors/okno/app/mod"
	"net/http"
)

// Create appends post path prefix for a database write
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
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
	JDB.Write(path+"/"+col, id, post)

}

var decoder = schema.NewDecoder()

// Change host of JDB
func Host(h string) *jdb {
	JDB.path = h
	return JDB
}
