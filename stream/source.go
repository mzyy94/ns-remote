package stream

import (
	"log"
	"math"
	"sync"

	"github.com/notedit/gst"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

// MediaSource is used to control pipelines
type MediaSource struct {
	videoPipeline *VideoPipeline
	audioPipeline *AudioPipeline
	videoChannel  chan struct{}
	audioChannel  chan struct{}
	waitGroup     sync.WaitGroup
	mutex         sync.Mutex
	IsLinked      bool
	streamer      *WebRTCStreamer
}

// NewMediaSource creates MediaSource from audio/video device names
func NewMediaSource(videosrc, audiosrc *string) *MediaSource {
	return &MediaSource{
		videoPipeline: NewVideoPipeline(videosrc),
		audioPipeline: NewAudioPipeline(audiosrc),
	}
}

// Link connects media source and stream outputs
func (p *MediaSource) Link(mediaStreamer *WebRTCStreamer) {
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
	p.streamer = mediaStreamer
}

// Unlink makes stop streaming
func (p *MediaSource) Unlink() {
	defer p.mutex.Unlock()
	p.mutex.Lock()
	if !p.IsLinked {
		return
	}

	close(p.videoChannel)
	close(p.audioChannel)
	p.waitGroup.Wait()

	log.Printf("* last stream state: %s\n", p.streamer.peerConnection.ConnectionState().String())

	if p.streamer.peerConnection.ConnectionState() != webrtc.PeerConnectionStateClosed {
		p.streamer.peerConnection.Close()
	}

	p.IsLinked = false
}

func startSampleTransfer(pipeline *gst.Pipeline, track *webrtc.Track, stop chan struct{}, waitGroup *sync.WaitGroup) {
	pipeline.SetState(gst.StatePlaying)
	sink := pipeline.GetByName("sink")
	waitGroup.Add(1)
	sampleChan := make(chan gst.Sample)

	log.Printf("-- start sample transfer of track %s\n", track.Kind().String())

	go func() {
		for {
			sample, err := sink.PullSample()
			if err != nil {
				log.Printf("Pull error on track %s: %s", track.Kind().String(), err.Error())
				return
			}
			sampleChan <- *sample
		}
	}()

	go func() {
		for {
			select {
			case <-stop:
				pipeline.SetState(gst.StateNull)
				waitGroup.Done()
				return
			case sample := <-sampleChan:
				samples := uint32(math.Round(float64(track.Codec().ClockRate) * (float64(sample.Duration) / 1000000000)))
				if err := track.WriteSample(media.Sample{Data: sample.Data, Samples: samples}); err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

// CheckGStreamerPlugins returns whether the GStreamer can be used
func CheckGStreamerPlugins() error {
	return gst.CheckPlugins([]string{
		"videotestsrc", "x264", "app",
		"audiotestsrc", "audioconvert", "audioresample", "opus",
	})
}
