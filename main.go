package main

import (
	"flag"
	"log"

	"ns-remote/server"
	"ns-remote/stream"

	"github.com/mzyy94/nscon"
)

func main() {
	var (
		videosrc = flag.String("video", "videotestsrc", "gstreamer video src")
		audiosrc = flag.String("audio", "audiotestsrc", "gstreamer audio src")
		device   = flag.String("device", "/dev/hidg0", "simulating hid gadget path")
		name     = flag.String("name", "procon", "configfs directory name")
	)
	flag.Parse()
	if err := stream.CheckGStreamerPlugins(); err != nil {
		log.Fatal(err)
	}
	mediaSource := stream.NewMediaSource(*videosrc, *audiosrc)
	controller := nscon.NewController(*device, *name)

	server.StartHTTPServer(mediaSource, controller)
}
