package main

import (
	"ns-remote/server"
	"ns-remote/stream"
)

func main() {
	mediaSource := stream.MediaSource{}
	mediaSource.Setup()

	server.StartHTTPServer(mediaSource)
}
