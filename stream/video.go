package stream

import (
	"github.com/notedit/gst"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

var pipeline *gst.Pipeline

// SetupVideoPipeline is..
func SetupVideoPipeline() {
	var err error
	pipeline, err = gst.PipelineNew("video-pipeline")

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

	sink, _ := gst.ElementFactoryMake("appsink", "sink")

	pipeline.AddMany(source, filter, convert, convertFilter, sink)

	source.Link(filter)
	filter.Link(convert)
	convert.Link(convertFilter)
	convertFilter.Link(sink)

	pipeline.SetState(gst.StatePaused)
}

// StartSampleTransfer is..
func StartSampleTransfer(track *webrtc.Track) {
	pipeline.SetState(gst.StatePlaying)
	sink := pipeline.GetByName("sink")

	for {
		sample, err := sink.PullSample()
		if err != nil {
			panic(err)
		}
		samples := uint32(90000 * (float32(sample.Duration) / 1000000000))
		if err := track.WriteSample(media.Sample{Data: sample.Data, Samples: samples}); err != nil {
			panic(err)
		}
	}
}
