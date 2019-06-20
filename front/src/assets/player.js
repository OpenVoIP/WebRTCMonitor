// h5splayer.js 不开源, 重新实现一个简单实现

function NewH5sPlayerRTC(conf) {
  this.url = conf.url
  this.videoid = conf.videoid
}

NewH5sPlayerRTC.prototype.disconnect = function () {
  console.log('断开连接')
}

NewH5sPlayerRTC.prototype.connect = function () {
  this.pc = new RTCPeerConnection()
  // let pc = new RTCPeerConnection({
  //   iceServers: [
  //     {
  //       urls: 'stun:stun.l.google.com:19302'
  //     }
  //   ]
  // })

  this.pc.ontrack = (event) => {
    let video = document.getElementById(this.videoid)
    video.srcObject = event.streams[0]
    video.autoplay = true
    video.controls = true
  }

  this.pc.oniceconnectionstatechange = e => {
    console.log(this.pc.iceConnectionState)
  }

  this.pc.onicecandidate = event => {
    if (event.candidate === null) {
      let localSessionDescription = btoa(this.pc.localDescription.sdp)
      $.post('http://localhost:8080/recive', { data: localSessionDescription }, (data) => {
        let remoteSessionDescription = data
        if (remoteSessionDescription === '') {
          return console.error('Session Description must not be empty')
        }

        try {
          this.pc.setRemoteDescription(new RTCSessionDescription({ type: 'answer', sdp: atob(remoteSessionDescription) }))
        } catch (e) {
          console.error(e)
        }
      })
    }
  }

  this.pc.createOffer({ offerToReceiveVideo: true, offerToReceiveAudio: true }).then(d => this.pc.setLocalDescription(d)).catch((e) => {
    console.error(e)
  })
}

export {NewH5sPlayerRTC}
