package main

import (
	"webrtc-monitor/rtsp"
	"webrtc-monitor/rtspserver"
	"webrtc-monitor/util"
	"webrtc-monitor/web"
)

// main 开始
func main() {

	//读取配置文件
	util.ReadConf()

	go web.StartHTTPServer()
	go rtsp.StartRTSPServer()
	go rtspserver.StartSterver()

	select {}
}
