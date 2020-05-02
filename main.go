package main

import (
	"fmt"

	"github.com/notedit/gst"
)

func main() {
	pipeline, err := gst.ParseLaunch("videotestsrc ! capsfilter caps=video/x-raw,width=1280,height=720 ! autovideosink")

	if err != nil {
		panic(err)
	}

	fmt.Println()

	pipeline.SetState(gst.StatePlaying)

	bus := pipeline.GetBus()

	for {
		message := bus.Pull(gst.MessageError | gst.MessageEos)
		fmt.Println("message:", message.GetName())
		if message.GetType() == gst.MessageEos {
			break
		}
	}
}
