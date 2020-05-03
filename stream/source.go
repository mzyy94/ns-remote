package stream

import (
	"log"
	"math"

	"github.com/notedit/gst"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

// MediaSource is..
type MediaSource struct {
	videoPipeline *VideoPipeline
	audioPipeline *AudioPipeline
	videoChannel  chan struct{}
	audioChannel  chan struct{}
	IsLinked      bool
}

// NewMediaSource is..
func NewMediaSource() (p MediaSource) {
	p.videoPipeline = NewVideoPipeline()
	p.audioPipeline = NewAudioPipeline()
	return
}

// Link is..
func (p *MediaSource) Link(mediaStreamer WebRTCStreamer) {
	if p.IsLinked {
		return
	}
	p.videoChannel = make(chan struct{})
	p.audioChannel = make(chan struct{})

	startSampleTransfer(p.videoPipeline, mediaStreamer.VideoTrack, p.videoChannel)
	startSampleTransfer(p.audioPipeline, mediaStreamer.AudioTrack, p.audioChannel)

	mediaStreamer.peerConnection.OnConnectionStateChange(func(connectionState webrtc.PeerConnectionState) {
		if connectionState == webrtc.PeerConnectionStateClosed {
			p.Unlink()
		}
	})
	p.IsLinked = true
}

// Unlink is..
func (p *MediaSource) Unlink() {
	if p.videoChannel != nil {
		close(p.videoChannel)
		p.videoChannel = nil
	}
	if p.audioChannel != nil {
		close(p.audioChannel)
		p.audioChannel = nil
	}
	p.IsLinked = false
}

func startSampleTransfer(pipeline *gst.Pipeline, track *webrtc.Track, stop chan struct{}) {
	pipeline.SetState(gst.StatePlaying)
	sink := pipeline.GetByName("sink")

	go func() {
		for {
			sample, err := sink.PullSample()
			if err != nil {
				panic(err)
			}
			select {
			case <-stop:
				pipeline.SetState(gst.StateNull)
				return
			default:
				samples := uint32(math.Round(float64(track.Codec().ClockRate) * (float64(sample.Duration) / 1000000000)))
				if err := track.WriteSample(media.Sample{Data: sample.Data, Samples: samples}); err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

// CheckGStreamerPlugins is..
func CheckGStreamerPlugins() error {
	return gst.CheckPlugins([]string{
		"videotestsrc", "x264", "app",
		"audiotestsrc", "audioconvert", "audioresample", "opus",
	})
}
