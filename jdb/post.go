package jdb

import (
	"encoding/json"
	"net/http"
)

var (
	out interface{}
)

// Renders JSON
func render(w http.ResponseWriter, p interface{}) {
	js, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
