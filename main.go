package main

import (
	"flag"
	"log"

	"ns-remote/server"
	"ns-remote/stream"
)

func main() {
	var (
		videosrc = flag.String("video", "videotestsrc", "gstreamer video src")
		audiosrc = flag.String("audio", "audiotestsrc", "gstreamer audio src")
	)
	flag.Parse()
	if err := stream.CheckGStreamerPlugins(); err != nil {
		log.Fatal(err)
	}
	mediaSource := stream.NewMediaSource(*videosrc, *audiosrc)

	server.StartHTTPServer(mediaSource)
}
