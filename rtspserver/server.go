package rtspserver

import (
	"fmt"

	"github.com/djwackey/dorsvr/rtspserver"
)

// StartSterver 启动服务器
func StartSterver() {
	server := rtspserver.New(nil)
	portNum := 8554
	err := server.Listen(portNum)
	if err != nil {
		fmt.Printf("Failed to bind port: %d\n", portNum)
		return
	}

	if !server.SetupTunnelingOverHTTP(80) ||
		!server.SetupTunnelingOverHTTP(8000) ||
		!server.SetupTunnelingOverHTTP(8080) {
		fmt.Printf("We use port %d for optional RTSP-over-HTTP tunneling, "+
			"or for HTTP live streaming (for indexed Transport Stream files only).\n", server.HTTPServerPortNum())
	} else {
		fmt.Println("(RTSP-over-HTTP tunneling is not available.)")
	}

	urlPrefix := server.RtspURLPrefix()
	fmt.Println("This server's URL: " + urlPrefix + "<filename>.")

	server.Start()
}
