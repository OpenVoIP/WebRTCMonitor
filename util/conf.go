package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

//Conf 配置
type Conf struct {
	HTTPConf struct {
		HTTPPort  int  `json:"nHTTPPort"`
		HTTPSPort int  `json:"nHTTPSPort"`
		Auth      bool `json:"bAuth"`
	} `json:"http"`

	RTSPConf struct {
		RTSPSink bool `json:"bRTSPSlink"`
		RTSPPort int  `json:"nRTSPPort"`
		SSLPort  int  `json:"SSLPort"`
		Auth     bool `json:"bAuth"`
	} `json:"rtsp"`

	WebRTCConf struct {
		WebRTCSink      bool   `json:"bWebRTCSink"`
		CloudMode       bool   `json:"bCloudMode"`
		RelatedPublicIP string `json:"strRelatedPublicIp"`
		PortRangMin     int    `json:"nPortRangeMin"`
		PortRangMax     int    `json:"nPortRangeMax"`
	} `json:"webrtc"`

	UserConf struct {
		TokenAuth     bool `json:"bTokenAuth"`
		AnonymousView bool `json:"bAnonymousView"`
		Users         []struct {
			User     string `json:"strUser"`
			Passwd   string `json:"strPasswd"`
			UserType string `json:"strUserType"`
		} `json:"user"`
	} `json:"user"`

	SourceConf struct {
		ConnectType       string `json:"nConnectType"`
		RTSPType          string `json:"nRTSPType"`
		EnablePreRecord   bool   `json:"bEnablePreRecord"`
		PreRecordLength   int    `json:"nPreRecordLength"`
		EnableIPPortCheck bool   `json:"bEnableIpPortCheck"`

		Src []struct {
			Name              string `json:"strName"`
			Token             string `json:"strToken"`
			Type              string `json:"nType"`
			URL               string `json:"strUrl"`
			User              string `json:"strUser"`
			Passwd            string `json:"strPasswd"`
			PasswdEncrypt     bool   `json:"bPasswdEncrypt"`
			EnableAudio       bool   `json:"bEnableAudio"`
			ConnectType       string `json:"nConnectType"`
			RTSPType          string `json:"nRTSPType"`
			SrcIPAddress      string `json:"strSrcIpAddress"`
			SrcPort           string `json:"strSrcPort"`
			ChannelNumber     int    `json:"nChannelNumber"`
			RTSPPlayback      bool   `json:"bRTSPPlayback"`
			RTSPPlaybackSpeed bool   `json:"bRTSPPlaybackSpeed"`
		} `json:"src"`
	} `json:"source"`

	DeviceConf struct {
		ConnectType     string `json:"nConnectType"`
		EnablePreRecord bool   `json:"bEnablePreRecord"`
		PreRecordLength int    `json:"nPreRecordLength"`

		Src []struct {
			Name          string `json:"strName"`
			Token         string `json:"strToken"`
			Type          string `json:"strType"`
			URL           string `json:"strUser"`
			PasswdEncrypt bool   `json:"bPasswdEncrypt"`
			DevIPAddress  string `json:"strDevIpAddress"`
			DevPort       string `json:"strDevPort"`
			EnableAudio   bool   `json:"bEnableAudio"`
		} `json:"src"`
	} `json:"device"`
}

//ConfData 服务器数据
var ConfData Conf

//ReadConf 阅读配置文件
func ReadConf() {
	workDir, _ := os.Getwd()
	jsonFile, err := os.Open(filepath.Join(workDir, "/h5ss.conf"))
	if err != nil {
		log.Error(err)
	}

	log.Info("Successfully Opened h5ss.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &ConfData)
	if err != nil {
		log.Error(err)
	}
}
