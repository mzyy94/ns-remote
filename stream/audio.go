package stream

import (
	"github.com/notedit/gst"
)

// AudioPipeline is..
type AudioPipeline = gst.Pipeline

// NewAudioPipeline is..
func NewAudioPipeline(device *string) *AudioPipeline {
	pipeline, _ := gst.PipelineNew("audio-pipeline")

	var source *gst.Element
	if device == nil {
		source, _ = gst.ElementFactoryMake("audiotestsrc", "source")
		source.SetObject("is-live", true)
	} else {
		source, _ = gst.ElementFactoryMake("alsasrc", "source")
		source.SetObject("device", *device)
	}

	convert, _ := gst.ElementFactoryMake("audioconvert", "convert")
	resample, _ := gst.ElementFactoryMake("audioresample", "resample")
	encoder, _ := gst.ElementFactoryMake("opusenc", "encoder")

	sink, _ := gst.ElementFactoryMake("appsink", "sink")

	pipeline.AddMany(source, convert, resample, encoder, sink)

	source.Link(convert)
	convert.Link(resample)
	resample.Link(encoder)
	encoder.Link(sink)

	pipeline.SetState(gst.StatePaused)
	return pipeline
}
