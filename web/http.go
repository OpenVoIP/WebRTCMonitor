package web

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"encoding/base64"

	"webrtc-monitor/util"
	v1 "webrtc-monitor/web/api/v1"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/pion/webrtc/v2"
	//ice "github.com/pions/webrtc/internal/ice"
)

var VideoTrack *webrtc.Track

func init() {

}

// StartHTTPServer 开启 HTTP
func StartHTTPServer() {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	// static dir
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	FileServer(r, "/", http.Dir(filesDir))

	// walk
	if err := chi.Walk(r, walkFunc); err != nil {
		log.Errorf("Logging err: %v\n", err.Error())
	}

	// rest
	r.HandleFunc("/recive", HTTPHome)
	r.Put("/ping", Ping)
	r.Mount("//api/v1", v1.Router())

	go func() {
		address := fmt.Sprintf(":%d", util.ConfData.HTTPConf.HTTPPort)
		log.Info(address)
		err := http.ListenAndServe(address, r)
		if err != nil {
			log.Error(err)
		}
	}()
	select {}
}

// walkFunc
func walkFunc(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	log.Infof("%s %s\n", method, route)
	return nil
}

// Ping returns pong
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// FileServer 静态文件
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// HTTPHome 首页处理
func HTTPHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data := r.FormValue("data")
	sd, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println(err)
		return
	}
	// webrtc.RegisterDefaultCodecs()
	//peerConnection, err := webrtc.New(webrtc.RTCConfiguration{
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})
	vp8Track, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeH264, rand.Uint32(), "video", "pion2")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = peerConnection.AddTrack(vp8Track)
	if err != nil {
		log.Println(err)
		return
	}
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  string(sd),
	}
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		log.Println(err)
		return
	}
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write([]byte(base64.StdEncoding.EncodeToString([]byte(answer.SDP))))
	VideoTrack = vp8Track
}
