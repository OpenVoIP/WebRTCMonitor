package web

import (
	"net/http"
	"webrtc-monitor/rtsp"
)

// NewPeerConnection 建立 peer
func NewPeerConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data := r.FormValue("data")
	name := r.FormValue("name")
	res, err := rtsp.SetVideoTrack(data, name)
	if err != nil {
		w.Write([]byte("error"))
		return
	}
	w.Write([]byte(res))
}
