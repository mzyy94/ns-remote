package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pion/webrtc/v2"
)

// WebRTCOfferHandler is..
func WebRTCOfferHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	var offer webrtc.SessionDescription
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		panic(err)
	}
	answer := CreateAnswerFromOffer(offer)
	json.NewEncoder(w).Encode(&answer)
}

// StartHTTPServer is..
func StartHTTPServer() {
	r := mux.NewRouter()

	r.HandleFunc("/connect", WebRTCOfferHandler).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go func() {
		srv.ListenAndServe()
	}()
}
