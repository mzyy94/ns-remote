package main

import (
	"ns-remote/server"
	"ns-remote/stream"
)

func main() {
	videoPipeline := stream.VideoPipeline{}
	audioPipeline := stream.AudioPipeline{}
	mStreamer := server.MediaStreamer{}

	videoPipeline.Setup()
	audioPipeline.Setup()
	mStreamer.Setup()
	go videoPipeline.StartSampleTransfer(mStreamer.VideoTrack)
	go audioPipeline.StartSampleTransfer(mStreamer.AudioTrack)
	server.StartHTTPServer(&mStreamer)
}
