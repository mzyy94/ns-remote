package main

import (
	"fmt"

	"github.com/notedit/gst"
)

func main() {
	pipeline, err := gst.PipelineNew("video-pipeline")

	if err != nil {
		panic(err)
	}

	source, _ := gst.ElementFactoryMake("videotestsrc", "source")

	filter, _ := gst.ElementFactoryMake("capsfilter", "filter")
	videoCap := gst.CapsFromString("video/x-raw,width=1280,height=720")
	filter.SetObject("caps", videoCap)

	sink, _ := gst.ElementFactoryMake("autovideosink", "sink")

	pipeline.AddMany(source, filter, sink)

	source.Link(filter)
	filter.Link(sink)

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
