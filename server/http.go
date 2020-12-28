package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mzyy94/nscon"
)

// StartHTTPServer starts HTTP server
func StartHTTPServer(con *nscon.Controller) {
	controller = con

	r := mux.NewRouter()

	r.HandleFunc("/controller", controllerHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Println("Connect http://localhost:8000")

	log.Fatal(srv.ListenAndServe())
}
