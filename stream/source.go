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
}

// NewMediaSource is..
func NewMediaSource() (p MediaSource) {
	p.videoPipeline = NewVideoPipeline()
	p.audioPipeline = NewAudioPipeline()
	return
}

// Link is..
func (p *MediaSource) Link(mediaStreamer WebRTCStreamer) {
	if p.videoChannel != nil || p.audioChannel != nil {
		log.Println("Already established")
		return
	}
	p.videoChannel = make(chan struct{})
	p.audioChannel = make(chan struct{})

	startSampleTransfer(p.videoPipeline, mediaStreamer.VideoTrack, p.videoChannel)
	startSampleTransfer(p.audioPipeline, mediaStreamer.AudioTrack, p.audioChannel)

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

func startSampleTransfer(pipeline *gst.Pipeline, track *webrtc.Track, ch chan struct{}) {
	pipeline.SetState(gst.StatePlaying)
	sink := pipeline.GetByName("sink")

	go func() {
		for {
			sample, err := sink.PullSample()
			if err != nil {
				panic(err)
			}
			select {
			case <-ch:
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
