package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/notedit/gst"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

func main() {
	pipeline, err := gst.PipelineNew("video-pipeline")

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

	convert, _ := gst.ElementFactoryMake("x264enc", "x264enc")

	convertFilter, _ := gst.ElementFactoryMake("capsfilter", "convertFilter")
	h264VideoCap := gst.CapsFromString("video/x-h264,width=1280,height=720,stream-format=byte-stream")
	convertFilter.SetObject("caps", h264VideoCap)

	sink, _ := gst.ElementFactoryMake("appsink", "sink")

	pipeline.AddMany(source, filter, convert, convertFilter, sink)

	source.Link(filter)
	filter.Link(convert)
	convert.Link(convertFilter)
	convertFilter.Link(sink)

	// WebRTC setup
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	// Create a video track
	videoTrack, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeH264, rand.Uint32(), "video", "video")
	if err != nil {
		panic(err)
	}
	_, err = peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}

	fmt.Println("waiting for signal")
	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}

	in, err := bufio.NewReader(os.Stdin).ReadString('\n')
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &offer)
	if err != nil {
		panic(err)
	}

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Output the answer in base64 so we can paste it in browser
	b, err = json.Marshal(answer)
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(b))

	pipeline.SetState(gst.StatePlaying)

	for {
		sample, err := sink.PullSample()
		if err != nil {
			panic(err)
		}
		samples := uint32(90000 * (float32(sample.Duration) / 1000000000))
		if err := videoTrack.WriteSample(media.Sample{Data: sample.Data, Samples: samples}); err != nil {
			panic(err)
		}
	}
}
