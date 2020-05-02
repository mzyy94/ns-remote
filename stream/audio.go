package stream

import (
	"math"

	"github.com/notedit/gst"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

// AudioPipeline is..
type AudioPipeline struct {
	pipeline *gst.Pipeline
}

// Setup is..
func (a *AudioPipeline) Setup() {
	var err error
	a.pipeline, err = gst.PipelineNew("audio-pipeline")

	if err != nil {
		panic(err)
	}

	source, _ := gst.ElementFactoryMake("audiotestsrc", "source")
	source.SetObject("is-live", true)

	convert, _ := gst.ElementFactoryMake("audioconvert", "convert")
	resample, _ := gst.ElementFactoryMake("audioresample", "resample")
	encoder, _ := gst.ElementFactoryMake("opusenc", "encoder")

	sink, _ := gst.ElementFactoryMake("appsink", "sink")

	a.pipeline.AddMany(source, convert, resample, encoder, sink)

	source.Link(convert)
	convert.Link(resample)
	resample.Link(encoder)
	encoder.Link(sink)

	a.pipeline.SetState(gst.StatePaused)
}

// StartSampleTransfer is..
func (a *AudioPipeline) StartSampleTransfer(track *webrtc.Track) {
	a.pipeline.SetState(gst.StatePlaying)
	sink := a.pipeline.GetByName("sink")

	for {
		sample, err := sink.PullSample()
		if err != nil {
			panic(err)
		}
		samples := uint32(math.Round(48000 * (float64(sample.Duration) / 1000000000)))
		if err := track.WriteSample(media.Sample{Data: sample.Data, Samples: samples}); err != nil {
			panic(err)
		}
	}
}
