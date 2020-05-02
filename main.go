package main

import (
	"ns-remote/server"
	"ns-remote/stream"
)

func main() {
	server.StartHTTPServer()
	stream.SetupVideoPipeline()
	videoTrack := server.SetupWebRTC()
	stream.StartSampleTransfer(videoTrack)
}
