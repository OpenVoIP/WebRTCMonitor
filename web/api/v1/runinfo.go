package v1

import (
	"net/http"

	"github.com/go-chi/render"
)

// RunInfo 运行信息
type RunInfo struct {
	StrCPU        string `json:"strCPU"`
	StrFreeSpace  string `json:"strFreeSpace"`
	StrMemory     string `json:"strMemory"`
	StrNetworkIn  string `json:"strNetworkIn"`
	StrNetworkOut string `json:"strNetworkOut"`
	StrRunTime    string `json:"strRunTime"`
	StrTotalSpace string `json:"strTotalSpace"`
}

// Render RunInfo 返回数据
func (info *RunInfo) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// GetRunInfo 获取运行信息
func GetRunInfo(w http.ResponseWriter, r *http.Request) {
	// "strRunTime": "0H 0MIN",
	// "strCPU": "3%",
	// "strMemory": "26%",
	// "strNetworkIn": "0Kbps",
	// "strNetworkOut": "0Kbps",
	// "strTotalSpace": "9323Mbytes",
	// "strFreeSpace": "541Mbytes"

	info := &RunInfo{
		StrRunTime:    "0H 0MIN",
		StrCPU:        "3%",
		StrMemory:     "26%",
		StrNetworkIn:  "0Kbps",
		StrNetworkOut: "0Kbps",
		StrTotalSpace: "9323Mbytes",
		StrFreeSpace:  "514Mbytes",
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, info)
}
