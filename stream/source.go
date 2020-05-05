package stream

import (
	"log"
	"math"
	"sync"
	"time"

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
	waitGroup     sync.WaitGroup
	mutex         sync.Mutex
	IsLinked      bool
}

// NewMediaSource is..
func NewMediaSource(videosrc, audiosrc *string) *MediaSource {
	return &MediaSource{
		videoPipeline: NewVideoPipeline(videosrc),
		audioPipeline: NewAudioPipeline(audiosrc),
	}
}

// Link is..
func (p *MediaSource) Link(mediaStreamer WebRTCStreamer) {
	defer p.mutex.Unlock()
	p.mutex.Lock()
	if p.IsLinked {
		return
	}
	p.videoChannel = make(chan struct{})
	p.audioChannel = make(chan struct{})

	startSampleTransfer(p.videoPipeline, mediaStreamer.VideoTrack, p.videoChannel, &p.waitGroup)
	startSampleTransfer(p.audioPipeline, mediaStreamer.AudioTrack, p.audioChannel, &p.waitGroup)

	mediaStreamer.peerConnection.OnConnectionStateChange(func(connectionState webrtc.PeerConnectionState) {
		if connectionState == webrtc.PeerConnectionStateClosed {
			p.Unlink()
		}
	})
	p.IsLinked = true
}

// Unlink is..
func (p *MediaSource) Unlink() {
	defer p.mutex.Unlock()
	p.mutex.Lock()
	if !p.IsLinked {
		return
	}

	close(p.videoChannel)
	close(p.audioChannel)
	p.waitGroup.Wait()
	p.IsLinked = false
}

func startSampleTransfer(pipeline *gst.Pipeline, track *webrtc.Track, stop chan struct{}, waitGroup *sync.WaitGroup) {
	pipeline.SetState(gst.StatePlaying)
	sink := pipeline.GetByName("sink")
	waitGroup.Add(1)

	bus := pipeline.GetBus()

	go func() {
		defer func() {
			pipeline.SetState(gst.StateNull)
			waitGroup.Done()
		}()

	Loop:
		for {
			if !bus.HavePending() {
				return
			}
			msg := bus.Pop()
			switch msg.GetType() {
			case gst.MessageStateChanged:
				break Loop
			case gst.MessageStreamStart:
				break Loop
			case gst.MessageError:
				pipeline.SetState(gst.StateNull)
				time.Sleep(3 * time.Second)
				pipeline.SetState(gst.StatePlaying)
			}
		}
		log.Printf("-- start sample transfer of track %s\n", track.Label())

		for {
			sample, err := sink.PullSample()
			if err != nil {
				panic(err)
			}
			select {
			case <-stop:
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
