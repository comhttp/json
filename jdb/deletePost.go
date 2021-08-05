package jdb

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// Delete  data from the database
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["host"]
	col := mux.Vars(r)["col"]
	id := mux.Vars(r)["slug"]
	DeletePost(path, col, id)
	return
}

// Delete  data from the database
func DeletePost(path, col, id string) {
	if err := JDB.Delete(path+"/"+col, id); err != nil {
		fmt.Println("Error", err)
	}
	return
}
