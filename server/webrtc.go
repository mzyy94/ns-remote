package server

import (
	"log"
	"math/rand"

	"github.com/pion/webrtc/v2"
)

// MediaStreamer is..
type MediaStreamer struct {
	peerConnection *webrtc.PeerConnection
	VideoTrack     *webrtc.Track
}

// Setup is ..
func (m *MediaStreamer) Setup() {
	// WebRTC setup
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	var err error
	m.peerConnection, err = webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	m.peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("Connection State has changed %s \n", connectionState.String())
	})

	// Create a video track
	m.VideoTrack, err = m.peerConnection.NewTrack(webrtc.DefaultPayloadTypeH264, rand.Uint32(), "video", "video")
	if err != nil {
		panic(err)
	}
	_, err = m.peerConnection.AddTrack(m.VideoTrack)
	if err != nil {
		panic(err)
	}
	return
}

// CreateAnswerFromOffer is..
func (m *MediaStreamer) CreateAnswerFromOffer(offer webrtc.SessionDescription) webrtc.SessionDescription {
	// Set the remote SessionDescription
	err := m.peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create an answer
	answer, err := m.peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = m.peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	return answer
}