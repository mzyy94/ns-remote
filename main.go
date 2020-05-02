package main

import (
	"ns-remote/server"
	"ns-remote/stream"
)

func main() {
	stream.SetupVideoPipeline()
	mStreamer := server.MediaStreamer{}
	mStreamer.Setup()
	go func() {
		stream.StartSampleTransfer(mStreamer.VideoTrack)
	}()
	server.StartHTTPServer(&mStreamer)
}
