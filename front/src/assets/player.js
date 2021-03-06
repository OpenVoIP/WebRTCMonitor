// h5splayer.js 不开源, 重新实现一个简单实现

let hostname = window.location.hostname;

function NewH5sPlayerRTC (conf) {
  this.url = conf.url
  this.videoid = conf.videoid
  this.name = conf.name
}

NewH5sPlayerRTC.prototype.disconnect = function () {
  console.log('断开连接')
}

NewH5sPlayerRTC.prototype.connect = function () {
  this.pc = new RTCPeerConnection({
    // iceServers: [
    //   {
    //     urls: 'stun:stun.l.google.com:19302'
    //   }
    // ]
  })

  this.pc.ontrack = (event) => {
    let video = document.getElementById(this.videoid)
    video.srcObject = event.streams[0]
    video.autoplay = true
    video.controls = false
  }

  this.pc.oniceconnectionstatechange = e => {
    console.log(this.pc.iceConnectionState)
  }

  this.pc.onicecandidate = event => {
    if (event.candidate === null) {
      let localSessionDescription = btoa(this.pc.localDescription.sdp)

      $.post(`http://${hostname}:8082/receive`, { data: localSessionDescription, name: this.name }, (data) => {
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
