package rtsp

import (
	"encoding/base64"
	"math/rand"
	"webrtc-monitor/util"

	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
	log "github.com/sirupsen/logrus"
)

//NameClientMap 缓存
var NameClientMap = make(map[string]*Client)

//StartRTSPServer 开启 RTSP 服务
func StartRTSPServer() {
	for _, src := range util.ConfData.SourceConf.Src {
		client := ClientNew()
		client.URL = src.URL
		client.Debug = false
		client.Name = src.Name
		NameClientMap[src.Name] = client

		sps := []byte{}
		pps := []byte{}
		fuBuffer := []byte{}
		count := 0

		syncCount := 0
		preTS := 0
		writeNALU := func(sync bool, ts int, payload []byte) {
			// if DataChanelTest != nil && preTS != 0 {
			// 	DataChanelTest <- webrtc.RTCSample{Data: payload, Samples: uint32(ts - preTS)}
			// }

			// 数据写入 web 端
			for _, track := range client.VideoTracks {
				if track != nil && preTS != 0 {
					s := ts - preTS
					go track.WriteSample(media.Sample{Data: payload, Samples: uint32(s)})
				}
			}
			preTS = ts

		}
		handleNALU := func(nalType byte, payload []byte, ts int64) {
			if nalType == 7 {
				if len(sps) == 0 {
					sps = payload
				}
				//	writeNALU(true, int(ts), payload)
			} else if nalType == 8 {
				if len(pps) == 0 {
					pps = payload
				}
				//	writeNALU(true, int(ts), payload)
			} else if nalType == 5 {
				syncCount++
				lastkeys := append([]byte("\000\000\001"+string(sps)+"\000\000\001"+string(pps)+"\000\000\001"), payload...)

				writeNALU(true, int(ts), lastkeys)
			} else {
				if syncCount > 0 {
					writeNALU(false, int(ts), payload)
				}
			}
		}

		if err := client.Open(); err != nil {
			log.Errorf("[RTSP] Error %+v", err)
			client.Close()
			return
		}

		for {
			select {
			case <-client.Signals:
				log.Error("Exit signals by rtsp")
				return
			case data := <-client.Outgoing:
				count += len(data)
				log.Debug("receive  rtp packet size", len(data), "receive all packet size", count)
				if data[0] == 36 && data[1] == 0 {
					cc := data[4] & 0xF
					rtphdr := 12 + cc*4
					ts := (int64(data[8]) << 24) + (int64(data[9]) << 16) + (int64(data[10]) << 8) + (int64(data[11]))
					packno := (int64(data[6]) << 8) + int64(data[7])
					if false {
						log.Info("packet num", packno)
					}
					nalType := data[4+rtphdr] & 0x1F
					if nalType >= 1 && nalType <= 23 {
						if nalType == 6 {
							continue
						}
						handleNALU(nalType, data[4+rtphdr:], ts)
					} else if nalType == 28 {
						isStart := data[4+rtphdr+1]&0x80 != 0
						isEnd := data[4+rtphdr+1]&0x40 != 0
						nalType := data[4+rtphdr+1] & 0x1F
						nal := data[4+rtphdr]&0xE0 | data[4+rtphdr+1]&0x1F
						if isStart {
							fuBuffer = []byte{0}
						}
						fuBuffer = append(fuBuffer, data[4+rtphdr+2:]...)
						if isEnd {
							fuBuffer[0] = nal
							handleNALU(nalType, fuBuffer, ts)
						}
					}
				} else if data[0] == 36 && data[1] == 2 {
					log.Info("data[0] == 36 data[1] == 2")
					//cc := data[4] & 0xF
					//rtphdr := 12 + cc*4
					//payload := data[4+rtphdr+4:]
				}
			}
		}
	}
}

//SetVideoTrack 设置track
func SetVideoTrack(data string, name string) (string, error) {
	sd, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Error(err)
		return "", err
	}
	// webrtc.RegisterDefaultCodecs()
	//peerConnection, err := webrtc.New(webrtc.RTCConfiguration{
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			// {
			// 	URLs: []string{"stun:stun.l.google.com:19302"},
			// },
		},
	})

	if err != nil {
		log.Panic(err)
	}
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Infof("ICE Connection State has changed %s \n", connectionState.String())
	})

	vp8Track, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeH264, rand.Uint32(), "video", "pion2")
	if err != nil {
		log.Error(err)
		return "", err
	}
	_, err = peerConnection.AddTrack(vp8Track)
	if err != nil {
		log.Error(err)
		return "", err
	}
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  string(sd),
	}
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		log.Error(err)
		return "", err
	}
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if client, exist := NameClientMap[name]; exist {
		client.VideoTracks = append(client.VideoTracks, vp8Track)
	}
	return base64.StdEncoding.EncodeToString([]byte(answer.SDP)), nil
}
