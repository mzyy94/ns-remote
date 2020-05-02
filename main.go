package main

import (
	"log"

	"ns-remote/server"
	"ns-remote/stream"
)

func main() {
	if err := stream.CheckGStreamerPlugins(); err != nil {
		log.Fatal(err)
	}
	mediaSource := stream.MediaSource{}
	mediaSource.Setup()

	server.StartHTTPServer(mediaSource)
}
