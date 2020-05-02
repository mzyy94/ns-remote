package stream

import (
	"math"

	"github.com/notedit/gst"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

// VideoPipeline is..
type VideoPipeline struct {
	pipeline *gst.Pipeline
}

// Setup is..
func (v *VideoPipeline) Setup() {
	var err error
	v.pipeline, err = gst.PipelineNew("video-pipeline")

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

	encoder, _ := gst.ElementFactoryMake("x264enc", "encoder")

	encodeFilter, _ := gst.ElementFactoryMake("capsfilter", "encodeFilter")
	h264VideoCap := gst.CapsFromString("video/x-h264,width=1280,height=720,stream-format=byte-stream,profile=constrained-baseline")
	encodeFilter.SetObject("caps", h264VideoCap)

	sink, _ := gst.ElementFactoryMake("appsink", "sink")

	v.pipeline.AddMany(source, filter, encoder, encodeFilter, sink)

	source.Link(filter)
	filter.Link(encoder)
	encoder.Link(encodeFilter)
	encodeFilter.Link(sink)

	v.pipeline.SetState(gst.StatePaused)
}

// StartSampleTransfer is..
func (v *VideoPipeline) StartSampleTransfer(track *webrtc.Track) {
	v.pipeline.SetState(gst.StatePlaying)
	sink := v.pipeline.GetByName("sink")

	for {
		sample, err := sink.PullSample()
		if err != nil {
			panic(err)
		}
		samples := uint32(math.Round(90000 * (float64(sample.Duration) / 1000000000)))
		if err := track.WriteSample(media.Sample{Data: sample.Data, Samples: samples}); err != nil {
			panic(err)
		}
	}
}
