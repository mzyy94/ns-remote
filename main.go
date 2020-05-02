package main

import (
	"ns-remote/server"
	"ns-remote/stream"
)

func main() {
	videoPipeline := stream.VideoPipeline{}
	mStreamer := server.MediaStreamer{}

	videoPipeline.Setup()
	mStreamer.Setup()
	go func() {
		videoPipeline.StartSampleTransfer(mStreamer.VideoTrack)
	}()
	server.StartHTTPServer(&mStreamer)
}
