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
		demo     = flag.Bool("demo", false, "use videotestsrc and audiotestsrc for A/V input")
		videosrc = flag.String("video", "/dev/video0", "v4l2 src device")
		audiosrc = flag.String("audio", "hw:0,0", "alsa src device")
		device   = flag.String("device", "/dev/hidg0", "simulating hid gadget path")
	)
	flag.Parse()
	if *demo {
		videosrc = nil
		audiosrc = nil
	}

	if err := stream.CheckGStreamerPlugins(); err != nil {
		log.Fatal(err)
	}

	mediaSource := stream.NewMediaSource(videosrc, audiosrc)
	controller := nscon.NewController(*device)

	server.StartHTTPServer(mediaSource, controller)
}
