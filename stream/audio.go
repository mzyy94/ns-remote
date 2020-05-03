package stream

import (
	"github.com/notedit/gst"
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
