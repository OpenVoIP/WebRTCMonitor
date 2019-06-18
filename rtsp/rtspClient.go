package rtsp

import (
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 视频宽高
var (
	VideoWidth  int
	VideoHeight int
)

// Client rtsp 客户端
type Client struct {
	socket   net.Conn
	outgoing chan []byte
	signals  chan bool
	host     string
	port     string
	uri      string
	auth     bool
	login    string
	password string
	session  string
	responce string
	bauth    string
	track    []string
	cseq     int
	videow   int
	videoh   int
}

// ClientNew init client
func ClientNew() *Client {
	Obj := &Client{
		cseq:     1,
		signals:  make(chan bool, 1),
		outgoing: make(chan []byte, 100000),
	}
	return Obj
}

//Client 处理数据
func (client *Client) Client(rtsp_url string) (bool, string) {
	if !client.ParseUrl(rtsp_url) {
		return false, "url error"
	}
	if !client.Connect() {
		return false, "connect error"
	}
	if !client.Write("OPTIONS " + client.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + "\r\n\r\n") {
		return false, "error OPTIONS"
	}
	if status, message := client.Read(); !status {
		return false, "connection lost"
	} else if status && strings.Contains(message, "Digest") {
		if !client.AuthDigest("OPTIONS", message) {
			return false, "Unautorized Digest"
		}
	} else if status && strings.Contains(message, "Basic") {
		if !client.AuthBasic("OPTIONS", message) {
			return false, "Unautorized Basic"
		}
	} else if !strings.Contains(message, "200") {
		return false, "err OPTIONS not status code 200 OK " + message
	}

	////////////PHASE 2 DESCRIBE
	log.Println("DESCRIBE " + client.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + client.bauth + "\r\n\r\n")
	if !client.Write("DESCRIBE " + client.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + client.bauth + "\r\n\r\n") {
		return false, "error DESCRIBE query"
	}
	if status, message := client.Read(); !status {
		return false, "DESCRIBE connection lost"
	} else if status && strings.Contains(message, "Digest") {
		if !client.AuthDigest("DESCRIBE", message) {
			return false, "Unautorized Digest"
		}
	} else if status && strings.Contains(message, "Basic") {
		if !client.AuthBasic("DESCRIBE", message) {
			return false, "Unautorized Basic"
		}
	} else if !strings.Contains(message, "200") {
		return false, "error DESCRIBE not status code 200 OK " + message
	} else {
		log.Println(message)
		client.track = client.ParseMedia(message)

	}
	if len(client.track) == 0 {
		return false, "error track not found "
	}
	//PHASE 3 SETUP
	log.Println("SETUP " + client.uri + "/" + client.track[0] + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + "\r\nTransport: RTP/AVP/TCP;unicast;interleaved=0-1" + client.bauth + "\r\n\r\n")
	if !client.Write("SETUP " + client.uri + "/" + client.track[0] + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + "\r\nTransport: RTP/AVP/TCP;unicast;interleaved=0-1" + client.bauth + "\r\n\r\n") {
		return false, ""
	}
	if status, message := client.Read(); !status {
		return false, "erro SETUP read"

	} else if !strings.Contains(message, "200") {
		if strings.Contains(message, "401") {
			str := client.AuthDigest_Only("SETUP", message)
			if !client.Write("SETUP " + client.uri + "/" + client.track[0] + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + "\r\nTransport: RTP/AVP/TCP;unicast;interleaved=0-1" + client.bauth + str + "\r\n\r\n") {
				return false, ""
			}
			if status, message := client.Read(); !status {
				return false, "error SETUP read"

			} else if !strings.Contains(message, "200") {

				return false, "error SETUP not status code 200 OK " + message

			} else {
				client.session = ParseSession(message)
			}
		} else {
			return false, "error SETUP not status code 200 OK " + message
		}
	} else {
		log.Println(message)
		client.session = ParseSession(message)
		log.Println(client.session)
	}
	if len(client.track) > 1 {

		if !client.Write("SETUP " + client.uri + "/" + client.track[1] + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + "\r\nTransport: RTP/AVP/TCP;unicast;interleaved=2-3" + "\r\nSession: " + client.session + client.bauth + "\r\n\r\n") {
			return false, ""
		}
		if status, message := client.Read(); !status {
			return false, "error SETUP Audio track"

		} else if !strings.Contains(message, "200") {
			if strings.Contains(message, "401") {
				str := client.AuthDigest_Only("SETUP", message)
				if !client.Write("SETUP " + client.uri + "/" + client.track[1] + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + "\r\nTransport: RTP/AVP/TCP;unicast;interleaved=2-3" + client.bauth + str + "\r\n\r\n") {
					return false, ""
				}
				if status, message := client.Read(); !status {
					return false, "error SETUP responce"

				} else if !strings.Contains(message, "200") {

					return false, "error SETUP not status code 200 OK " + message

				} else {
					log.Println(message)
					client.session = ParseSession(message)
				}
			} else {
				return false, "error SETUP not status code 200 OK " + message
			}
		} else {
			log.Println(message)
			client.session = ParseSession(message)
		}
	}

	log.Println("PLAY " + client.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + "\r\nSession: " + client.session + client.bauth + "\r\n\r\n")
	if !client.Write("PLAY " + client.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + "\r\nSession: " + client.session + client.bauth + "\r\n\r\n") {
		return false, ""
	}
	if status, message := client.Read(); !status {
		return false, "error PLAY connection lost"

	} else if !strings.Contains(message, "200") {
		if strings.Contains(message, "401") {
			str := client.AuthDigest_Only("PLAY", message)
			if !client.Write("PLAY " + client.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + "\r\nSession: " + client.session + client.bauth + str + "\r\n\r\n") {
				return false, ""
			}
			if status, message := client.Read(); !status {
				return false, "error PLAY connection lost"

			} else if !strings.Contains(message, "200") {

				return false, "error PLAY not status code 200 OK " + message

			} else {
				log.Print(message)
				go client.RtspRtpLoop()
				return true, "ok"
			}
		} else {
			return false, "error PLAY not status code 200 OK " + message
		}
	} else {
		log.Print(message)
		go client.RtspRtpLoop()
		return true, "ok"
	}
}

/*
	The RTP header has the following format:

    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |V=2|P|X|  CC   |M|     PT      |       sequence number         |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                           timestamp                           |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |           synchronization source (SSRC) identifier            |
   +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
   |            contributing source (CSRC) identifiers             |
   |                             ....                              |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   version (V): 2 bits
      client field identifies the version of RTP.  The version defined by
      client specification is two (2).  (The value 1 is used by the first
      draft version of RTP and the value 0 is used by the protocol
      initially implemented in the "vat" audio tool.)

   padding (P): 1 bit
      If the padding bit is set, the packet contains one or more
      additional padding octets at the end which are not part of the
      payload.  The last octet of the padding contains a count of how
      many padding octets should be ignored, including itself.  Padding
      may be needed by some encryption algorithms with fixed block sizes
      or for carrying several RTP packets in a lower-layer protocol data
      unit.

   extension (X): 1 bit
      If the extension bit is set, the fixed header MUST be followed by
      exactly one header extension, with a format defined in Section
      5.3.1.

*/
// RtspRtpLoop loop
func (client *Client) RtspRtpLoop() {
	defer func() {
		client.signals <- true
	}()
	header := make([]byte, 4)
	payload := make([]byte, 4096)
	syncB := make([]byte, 1)
	timer := time.Now()
	for {
		if int(time.Now().Sub(timer).Seconds()) > 50 {
			if !client.Write("OPTIONS " + client.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + "\r\nSession: " + client.session + client.bauth + "\r\n\r\n") {
				return
			}
			timer = time.Now()
		}
		client.socket.SetDeadline(time.Now().Add(50 * time.Second))
		if n, err := io.ReadFull(client.socket, header); err != nil || n != 4 {
			return
		}
		if header[0] != 36 {
			for {
				if n, err := io.ReadFull(client.socket, syncB); err != nil && n != 1 {
					return
				} else if syncB[0] == 36 {
					header[0] = 36
					if n, err := io.ReadFull(client.socket, header[1:]); err != nil && n == 3 {
						return
					}
					break
				}
			}
		}
		payloadLen := (int)(header[2])<<8 + (int)(header[3])
		if payloadLen > 4096 || payloadLen < 12 {
			log.Println("desync", client.uri, payloadLen)
			return
		}
		if n, err := io.ReadFull(client.socket, payload[:payloadLen]); err != nil || n != payloadLen {
			return
		} else {
			client.outgoing <- append(header, payload[:n]...)
		}
	}

}

// SendBufer buffer
func (client *Client) SendBufer(bufer []byte) {
	payload := make([]byte, 4096)
	for {
		if len(bufer) < 4 {
			log.Fatal("bufer small")
		}
		dataLength := (int)(bufer[2])<<8 + (int)(bufer[3])
		if dataLength > len(bufer)+4 {
			if n, err := io.ReadFull(client.socket, payload[:dataLength-len(bufer)+4]); err != nil {
				return
			} else {
				client.outgoing <- append(bufer, payload[:n]...)
				return
			}

		} else {
			client.outgoing <- bufer[:dataLength+4]
			bufer = bufer[dataLength+4:]
		}
	}
}

// Connect 连接
func (client *Client) Connect() bool {
	d := &net.Dialer{Timeout: 3 * time.Second}
	conn, err := d.Dial("tcp", client.host+":"+client.port)
	if err != nil {
		return false
	}
	client.socket = conn
	return true
}

// Write write
func (client *Client) Write(message string) bool {
	client.cseq += 1
	if _, e := client.socket.Write([]byte(message)); e != nil {
		return false
	}
	return true
}

// Read read
func (client *Client) Read() (bool, string) {
	buffer := make([]byte, 4096)
	if nb, err := client.socket.Read(buffer); err != nil || nb <= 0 {
		log.Println("socket read failed", err)
		return false, ""
	} else {
		return true, string(buffer[:nb])
	}
}

// AuthBasic 认证
func (client *Client) AuthBasic(phase string, message string) bool {
	client.bauth = "\r\nAuthorization: Basic " + b64.StdEncoding.EncodeToString([]byte(client.login+":"+client.password))
	if !client.Write(phase + " " + client.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + client.bauth + "\r\n\r\n") {
		return false
	}
	if status, message := client.Read(); status && strings.Contains(message, "200") {
		client.track = ParseMedia(message)
		return true
	}
	return false
}

// AuthDigest digest
func (client *Client) AuthDigest(phase string, message string) bool {
	nonce := ParseDirective(message, "nonce")
	realm := ParseDirective(message, "realm")
	hs1 := GetMD5Hash(client.login + ":" + realm + ":" + client.password)
	hs2 := GetMD5Hash(phase + ":" + client.uri)
	responce := GetMD5Hash(hs1 + ":" + nonce + ":" + hs2)
	dauth := "\r\n" + `Authorization: Digest username="` + client.login + `", realm="` + realm + `", nonce="` + nonce + `", uri="` + client.uri + `", response="` + responce + `"`
	if !client.Write(phase + " " + client.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(client.cseq) + dauth + "\r\n\r\n") {
		return false
	}
	if status, message := client.Read(); status && strings.Contains(message, "200") {
		client.track = ParseMedia(message)
		return true
	}
	return false
}

//AuthDigest_Only only
func (client *Client) AuthDigest_Only(phase string, message string) string {
	nonce := ParseDirective(message, "nonce")
	realm := ParseDirective(message, "realm")
	hs1 := GetMD5Hash(client.login + ":" + realm + ":" + client.password)
	hs2 := GetMD5Hash(phase + ":" + client.uri)
	responce := GetMD5Hash(hs1 + ":" + nonce + ":" + hs2)
	dauth := "\r\n" + `Authorization: Digest username="` + client.login + `", realm="` + realm + `", nonce="` + nonce + `", uri="` + client.uri + `", response="` + responce + `"`
	return dauth
}

// ParseUrl parse
func (client *Client) ParseUrl(rtsp_url string) bool {

	u, err := url.Parse(rtsp_url)
	if err != nil {
		return false
	}
	phost := strings.Split(u.Host, ":")
	client.host = phost[0]
	if len(phost) == 2 {
		client.port = phost[1]
	} else {
		client.port = "554"
	}
	client.login = u.User.Username()
	client.password, client.auth = u.User.Password()
	if u.RawQuery != "" {
		client.uri = "rtsp://" + client.host + ":" + client.port + u.Path + "?" + string(u.RawQuery)
	} else {
		client.uri = "rtsp://" + client.host + ":" + client.port + u.Path
	}
	return true
}

// Close close
func (client *Client) Close() {
	if client.socket != nil {
		client.socket.Close()
	}
}

// ParseDirective directive
func ParseDirective(header, name string) string {
	index := strings.Index(header, name)
	if index == -1 {
		return ""
	}
	start := 1 + index + strings.Index(header[index:], `"`)
	end := start + strings.Index(header[start:], `"`)
	return strings.TrimSpace(header[start:end])
}

// ParseSession session
func ParseSession(header string) string {
	mparsed := strings.Split(header, "\r\n")
	for _, element := range mparsed {
		if strings.Contains(element, "Session:") {
			if strings.Contains(element, ";") {
				fist := strings.Split(element, ";")[0]
				return fist[9:]
			} else {
				return element[9:]
			}
		}
	}
	return ""
}

// ParseMedia media
func ParseMedia(header string) []string {
	letters := []string{}
	mparsed := strings.Split(header, "\r\n")
	paste := ""

	if true {
		log.Println("headers", header)
	}

	for _, element := range mparsed {
		if strings.Contains(element, "a=control:") && !strings.Contains(element, "*") && strings.Contains(element, "tra") {
			paste = element[10:]
			if strings.Contains(element, "/") {
				striped := strings.Split(element, "/")
				paste = striped[len(striped)-1]
			}
			letters = append(letters, paste)
		}

		dimensionsPrefix := "a=x-dimensions:"
		if strings.HasPrefix(element, dimensionsPrefix) {
			dims := []int{}
			for _, s := range strings.Split(element[len(dimensionsPrefix):], ",") {
				v := 0
				fmt.Sscanf(s, "%d", &v)
				if v <= 0 {
					break
				}
				dims = append(dims, v)
			}
			if len(dims) == 2 {
				VideoWidth = dims[0]
				VideoHeight = dims[1]
			}
		}
	}
	return letters
}

//GetMD5Hash md5
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

//ParseMedia media
func (client *Client) ParseMedia(header string) []string {
	letters := []string{}
	mparsed := strings.Split(header, "\r\n")
	paste := ""
	for _, element := range mparsed {
		if strings.Contains(element, "a=control:") && !strings.Contains(element, "*") && strings.Contains(element, "tra") {
			paste = element[10:]
			if strings.Contains(element, "/") {
				striped := strings.Split(element, "/")
				paste = striped[len(striped)-1]
			}
			letters = append(letters, paste)
		}

		dimensionsPrefix := "a=x-dimensions:"
		if strings.HasPrefix(element, dimensionsPrefix) {
			dims := []int{}
			for _, s := range strings.Split(element[len(dimensionsPrefix):], ",") {
				v := 0
				fmt.Sscanf(s, "%d", &v)
				if v <= 0 {
					break
				}
				dims = append(dims, v)
			}
			if len(dims) == 2 {
				client.videow = dims[0]
				client.videoh = dims[1]
			}
		}
	}
	return letters
}
