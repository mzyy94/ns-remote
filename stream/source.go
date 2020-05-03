package stream

import (
	"log"

	"github.com/notedit/gst"
	"github.com/pion/webrtc/v2"
)

// MediaSource is..
type MediaSource struct {
	videoPipeline *VideoPipeline
	audioPipeline *AudioPipeline
	videoChannel  chan struct{}
	audioChannel  chan struct{}
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
	if p.videoChannel != nil || p.audioChannel != nil {
		log.Println("Already established")
		return
	}
	p.videoChannel = make(chan struct{})
	p.audioChannel = make(chan struct{})

	p.videoPipeline.StartSampleTransfer(mediaStreamer.VideoTrack, p.videoChannel)
	p.audioPipeline.StartSampleTransfer(mediaStreamer.AudioTrack, p.audioChannel)

	mediaStreamer.peerConnection.OnConnectionStateChange(func(connectionState webrtc.PeerConnectionState) {
		if connectionState == webrtc.PeerConnectionStateClosed {
			p.Stop()
		}
	})
}

// Stop is..
func (p *MediaSource) Stop() {
	if p.videoChannel == nil || p.audioChannel == nil {
		log.Println("Connection not established")
		return
	}
	close(p.videoChannel)
	close(p.audioChannel)
	p.videoChannel = nil
	p.audioChannel = nil
}

// CheckGStreamerPlugins is..
func CheckGStreamerPlugins() error {
	return gst.CheckPlugins([]string{
		"videotestsrc", "x264", "app",
		"audiotestsrc", "audioconvert", "audioresample", "opus",
	})
}
