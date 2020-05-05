package stream

import (
	"github.com/notedit/gst"
)

// VideoPipeline is..
type VideoPipeline = gst.Pipeline

// NewVideoPipeline is..
func NewVideoPipeline(videosrc string) *VideoPipeline {
	pipeline, _ := gst.PipelineNew("video-pipeline")

	source, _ := gst.ElementFactoryMake(videosrc, "source")
	if videosrc == "videotestsrc" {
		source.SetObject("is-live", true)
		source.SetObject("pattern", 18)
		source.SetObject("animation-mode", 1)
		source.SetObject("motion", 1)
	}

	filter, _ := gst.ElementFactoryMake("capsfilter", "filter")
	rawVideoCap := gst.CapsFromString("video/x-raw,width=1280,height=720")
	filter.SetObject("caps", rawVideoCap)

	var encoder *gst.Element
	if gst.CheckPlugins([]string{"video4linux2"}) == nil {
		var err error
		if encoder, err = gst.ElementFactoryMake("v4l2h264enc", "encoder"); err == nil {
			controls := gst.NewStructure("encode")
			controls.SetValue("h264_profile", 1)
			controls.SetValue("h264_level", 12)
			encoder.SetObject("extra-controls", controls)
		} else {
			encoder, _ = gst.ElementFactoryMake("x264enc", "encoder")
		}
	} else {
		encoder, _ = gst.ElementFactoryMake("x264enc", "encoder")
	}

	encodeFilter, _ := gst.ElementFactoryMake("capsfilter", "encodeFilter")
	h264VideoCap := gst.CapsFromString("video/x-h264,width=1280,height=720,stream-format=byte-stream,profile=constrained-baseline")
	encodeFilter.SetObject("caps", h264VideoCap)

	parser, _ := gst.ElementFactoryMake("h264parse", "parser")
	parser.SetObject("config-interval", -1)

	sink, _ := gst.ElementFactoryMake("appsink", "sink")

	pipeline.AddMany(source, filter, encoder, encodeFilter, parser, sink)

	source.Link(filter)
	filter.Link(encoder)
	encoder.Link(encodeFilter)
	encodeFilter.Link(parser)
	parser.Link(sink)

	pipeline.SetState(gst.StatePaused)
	return pipeline
}
