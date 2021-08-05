package main

import (
	"fmt"
	"net/http"
)

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}

type interceptResponseWriter struct {
	http.ResponseWriter
	errH func(http.ResponseWriter, int)
}

func (w *interceptResponseWriter) WriteHeader(status int) {
	if status >= http.StatusBadRequest {
		w.errH(w.ResponseWriter, status)
		w.errH = nil
	} else {
		w.ResponseWriter.WriteHeader(status)
	}
}

type ErrorHandler func(http.ResponseWriter, int)

func (w *interceptResponseWriter) Write(p []byte) (n int, err error) {
	if w.errH == nil {
		return len(p), nil
	}
	return w.ResponseWriter.Write(p)
}

func DefaultErrorHandler(w http.ResponseWriter, status int) {
	//t := template.Must(template.ParseFiles("errors/error.html"))
	//w.Header().Set("Content-Type", "text/html")
	//t.Execute(w, map[string]interface{}{"status": status})
	w.Header().Set("Content-Type", "text/html")
	//tpl.TemplateHandler("/templates/sys").ExecuteTemplate(w, "error_gohtml", map[string]interface{}{"status": status})

}

func InterceptHandler(next http.Handler, errH ErrorHandler) http.Handler {
	if errH == nil {
		errH = DefaultErrorHandler
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&interceptResponseWriter{w, errH}, r)
	})
}

func ErrorLog(err error) {
	if err != nil {
		fmt.Println("Error", err)
	}
}
