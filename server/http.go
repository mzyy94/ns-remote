package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ns-remote/stream"

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

	videoPipeline := stream.VideoPipeline{}
	audioPipeline := stream.AudioPipeline{}
	mStreamer := MediaStreamer{}

	videoPipeline.Setup()
	audioPipeline.Setup()
	mStreamer.Setup(offer)
	go videoPipeline.StartSampleTransfer(mStreamer.VideoTrack)
	go audioPipeline.StartSampleTransfer(mStreamer.AudioTrack)

	answer := mStreamer.CreateAnswerFromOffer(offer)
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
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
