package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/oknors/okno/app/mod"
	"github.com/oknors/okno/pkg/utl"
	"net/http"
)

// // Response Handler
func (j *API) ReadPostHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["host"]
	col := mux.Vars(r)["col"]
	id := mux.Vars(r)["slug"]
	render(w, ReadPost(j, path, col, id))
}

// // ReadData appends 'data' path prefix for a database read
func ReadPost(j *API, path, col, id string) mod.Post {
	data := mod.Post{}
	err := j.Read(path+"/"+col, id, &data)
	utl.ErrorLog(err)
	return data
}

// Rresponse Handler.
func (j *API) ReadPostCollectionHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["host"]
	col := mux.Vars(r)["col"]
	render(w, ReadPostCollection(j, path, col))
}

// Read all items from the database, unmarshaling the response.
func ReadPostCollection(j *API, path, col string) []mod.Post {
	var posts []mod.Post
	data, err := j.ReadAll(path + "/" + col)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, postInterface := range data {
		var p mod.Post
		if err := json.Unmarshal([]byte(postInterface), &p); err != nil {
			fmt.Println("Error", err)
		}
		posts = append(posts, p)
	}
	return posts
}
