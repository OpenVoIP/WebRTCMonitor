package v1

import (
	"net/http"
	"webrtc-monitor/util"

	"github.com/go-chi/render"
)

//SrcInfo 视频源信息
type SrcInfo struct {
	EnableAudio       bool   `json:"bEnableAudio"`
	Online            bool   `json:"bOnLine"`
	OnvifProfileAuto  bool   `json:"bOnvifProfileAuto"`
	PasswdEncrypt     bool   `json:"bPasswdEncrypt"`
	RTSPPlayback      bool   `json:"bRTSPPlayback"`
	Rec               int    `json:"nRec"`
	ChannelNumber     string `json:"nChannelNumber"`
	ConnectType       string `json:"nConnectType"`
	OriginalType      string `json:"nOriginalType"`
	RTSPPlaybackSpeed string `json:"nRTSPPlaybackSpeed"`
	RTSPType          string `json:"nRTSPType"`
	Type              string `json:"nType"`
	Name              string `json:"strName"`
	OnvifAddr         string `json:"strOnvifAddr"`
	OnvifProfileMain  string `json:"strOnvifProfileMain"`
	OnvifProfileSub   string `json:"strOnvifProfileSub"`
	OriginalToken     string `json:"strOriginalToken"`
	Passwd            string `json:"strPasswd"`
	PushURL           string `json:"strPushUrl"`
	ServerToken       string `json:"strServerToken"`
	ServerUuid        string `json:"strServerUuid"`
	SnapshotUrl       string `json:"strSnapshotUrl"`
	SrcIPAddress      string `json:"strSrcIpAddress"`
	SrcPort           string `json:"strSrcPort"`
	Token             string `json:"strToken"`
	URL               string `json:"strUrl"`
	User              string `json:"strUser"`
}

// Render RunInfo 返回数据
func (info *SrcInfo) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//GetSrcInfo 获取视频源
func GetSrcInfo(w http.ResponseWriter, r *http.Request) {
	infos := []*SrcInfo{}
	for _, src := range util.ConfData.SourceConf.Src {
		infos = append(infos, &SrcInfo{
			URL:    src.URL,
			Online: true,
			Token:  src.Token,
			Name:   src.Name,
		})
	}

	render.Status(r, http.StatusOK)
	render.DefaultResponder(w, r, render.M{"src": infos})
}
