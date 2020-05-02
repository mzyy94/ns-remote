package stream

import "github.com/notedit/gst"

// MediaSource is..
type MediaSource struct {
	videoPipeline *VideoPipeline
	audioPipeline *AudioPipeline
}

// Setup is..
func (p *MediaSource) Setup() {
	p.videoPipeline = new(VideoPipeline)
	p.audioPipeline = new(AudioPipeline)

	p.videoPipeline.Setup()
	p.audioPipeline.Setup()
}

// Link is..
func (p *MediaSource) Link(mediaStreamer WebRTCStreamer) {
	go p.videoPipeline.StartSampleTransfer(mediaStreamer.VideoTrack)
	go p.audioPipeline.StartSampleTransfer(mediaStreamer.AudioTrack)
}

// CheckGStreamerPlugins is..
func CheckGStreamerPlugins() error {
	return gst.CheckPlugins([]string{
		"videotestsrc", "x264", "app",
		"audiotestsrc", "audioconvert", "audioresample", "opus",
	})
}
