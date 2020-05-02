package main

import (
	"ns-remote/server"
	"ns-remote/stream"
)

func main() {
	stream.SetupVideoPipeline()
	videoTrack := server.SetupWebRTC()
	stream.StartSampleTransfer(videoTrack)
}
