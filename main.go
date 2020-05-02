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
	go func() {
		videoPipeline.StartSampleTransfer(mStreamer.VideoTrack)
	}()
	go func() {
		audioPipeline.StartSampleTransfer(mStreamer.AudioTrack)
	}()
	server.StartHTTPServer(&mStreamer)
}
