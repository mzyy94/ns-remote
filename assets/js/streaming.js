const pc = new RTCPeerConnection({
  iceServers: [
    {
      urls: "stun:stun.l.google.com:19302"
    }
  ]
});

pc.ontrack = event => {
  const video = document.querySelector("video");
  if (video.srcObject == null) {
    video.srcObject = event.streams[0];
  } else {
    video.srcObject.addTrack(event.track);
  }
};

pc.oniceconnectionstatechange = () => console.log(pc.iceConnectionState);
pc.onicecandidate = event => {
  if (event.candidate == null) {
    fetch("/connect", {
      method: "POST",
      body: JSON.stringify(pc.localDescription)
    })
      .then(res => Promise.all([res.json(), res.ok]))
      .then(([answer, ok]) => {
        if (!ok) {
          return Promise.reject(answer);
        }
        try {
          pc.setRemoteDescription(new RTCSessionDescription(answer));
        } catch (e) {
          Promise.reject(e);
        }
      })
      .catch(console.error);
  }
};

pc.addTransceiver("video", { direction: "recvonly" });
pc.addTransceiver("audio", { direction: "recvonly" });

pc.createOffer()
  .then(d => pc.setLocalDescription(d))
  .catch(console.error);

window.addEventListener("pointerdown", () => {
  document.querySelector("video").muted = false;
}, {once: true});
