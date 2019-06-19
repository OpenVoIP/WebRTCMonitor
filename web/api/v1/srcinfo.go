package v1

import (
	"net/http"

	"github.com/go-chi/render"
)

//SrcInfo 视频源信息
type SrcInfo struct {
	BEnableAudio        string `json:"bEnableAudio"`
	BOnline             string `json:"bOnLine"`
	BOnvifProfileAuto   string `json:"bOnvifProfileAuto"`
	BPasswdEncrypt      string `json:"bPasswdEncrypt"`
	BRTSPPlayback       string `json:"bRTSPPlayback"`
	NRec                string `json:"nRec"`
	NChannelNumber      string `json:"nChannelNumber"`
	NConnectType        string `json:"nConnectType"`
	NOriginalType       string `json:"nOriginalType"`
	NRTSPPlaybackSpeed  string `json:"nRTSPPlaybackSpeed"`
	NRTSPType           string `json:"nRTSPType"`
	NType               string `json:"nType"`
	StrName             string `json:" "strName"`
	StrOnvifAddr        string `json:" "strOnvifAddr"`
	StrOnvifProfileMain string `json:"strOnvifProfileMain"`
	StrOnvifProfileSub  string `json:"strOnvifProfileSub"`
	StrOriginalToken    string `json:"strOriginalToken"`
	StrPasswd           string `json:"strPasswd"`
	StrPushUrl          string `json:"strPushUrl"`
	StrServerToken      string `json:"strServerToken"`
	StrServerUuid       string `json:"strServerUuid"`
	StrSnapshotUrl      string `json:"strSnapshotUrl"`
	StrSrcIpAddress     string `json:"strSrcIpAddress"`
	StrSrcPort          string `json:"strSrcPort"`
	StrToken            string `json:"strToken"`
	StrUrl              string `json:"strUrl"`
	StrUser             string `json: "strUser"`
}

// Render RunInfo 返回数据
func (info *SrcInfo) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//GetSrcInfo 获取视频源
func GetSrcInfo(w http.ResponseWriter, r *http.Request) {
	infos := []*SrcInfo{
		{
			StrUrl: "rstp://admin:admin@192.168.11.63/live",
		},
	}
	render.Status(r, http.StatusOK)
	render.DefaultResponder(w, r, render.M{"src": infos})
}
