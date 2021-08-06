package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// Delete  data from the database
func (j *API) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["host"]
	col := mux.Vars(r)["col"]
	id := mux.Vars(r)["slug"]
	DeletePost(j, path, col, id)
	return
}

// Delete  data from the database
func DeletePost(j *API, path, col, id string) {
	if err := j.Delete(path+"/"+col, id); err != nil {
		fmt.Println("Error", err)
	}
	return
}
