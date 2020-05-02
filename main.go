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
	source.SetObject("is-live", true)
	source.SetObject("pattern", 18)
	source.SetObject("animation-mode", 1)
	source.SetObject("motion", 1)

	filter, _ := gst.ElementFactoryMake("capsfilter", "filter")
	rawVideoCap := gst.CapsFromString("video/x-raw,width=1280,height=720")
	filter.SetObject("caps", rawVideoCap)

	convert, _ := gst.ElementFactoryMake("x264enc", "x264enc")

	convertFilter, _ := gst.ElementFactoryMake("capsfilter", "convertFilter")
	h264VideoCap := gst.CapsFromString("video/x-h264,width=1280,height=720,stream-format=byte-stream")
	convertFilter.SetObject("caps", h264VideoCap)

	parser, _ := gst.ElementFactoryMake("h264parse", "parser")

	rtph264pay, _ := gst.ElementFactoryMake("rtph264pay", "rtph264pay")
	rtph264pay.SetObject("config-interval", -1)
	rtph264pay.SetObject("pt", 96)

	sink, _ := gst.ElementFactoryMake("appsink", "sink")

	pipeline.AddMany(source, filter, convert, convertFilter, parser, rtph264pay, sink)

	source.Link(filter)
	filter.Link(convert)
	convert.Link(convertFilter)
	convertFilter.Link(parser)
	parser.Link(rtph264pay)
	rtph264pay.Link(sink)

	pipeline.SetState(gst.StatePlaying)

	for {
		sample, err := sink.PullSample()
		if err != nil {
			panic(err)
		}
		fmt.Println("got sample", sample.Duration)
	}
}
