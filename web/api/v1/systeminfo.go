package v1

import (
	"net/http"

	"github.com/go-chi/render"
)

//SystemInfo 系统信息
type SystemInfo struct {
	StrVersion      string `json:"strVersion"`
	StrLicenseType  string `json:"strLicenseType"`
	StrHostID       string `json:"strHostId"`
	StrLicenseFull  string `json:"strLicenseFull"`
	StrChannelLimit string `json:"strChannelLimit"`
	StrEndtime      string `json:"strEndtime"`
}

// Render SystemInfo 返回数据
func (info *SystemInfo) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//GetSystemInfo 获取系统信息
func GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	// "strVersion": "r9.1.0618.19",
	// "strHostId": "OWI2MDUzNmVjMzIyOWE0ZjFhOTk4YjBkNDhlYzBiNzk=",
	// "strLicenseType": "None",
	// "strLicenseFull": "******",
	// "strChannelLimit": "unlimited",
	// "strEndtime": "unlimited"

	info := &SystemInfo{
		StrVersion:      "r9.1.0618.19",
		StrHostID:       "OWI2MDUzNmVjMzIyOWE0ZjFhOTk4YjBkNDhlYzBiNzk=",
		StrLicenseType:  "None",
		StrLicenseFull:  "******",
		StrChannelLimit: "unlimited",
		StrEndtime:      "unlimited",
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, info)
}
