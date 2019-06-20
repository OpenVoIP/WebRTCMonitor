package v1

import (
	"net/http"

	"github.com/go-chi/render"
)

//Snapshot 截屏
func Snapshot(w http.ResponseWriter, r *http.Request) {
	// token := r.URL.Query()["token"][0]
	// ivfFile, err := ivfwriter.New("output.ivf")
	// if err != nil {
	// 	panic(err)
	// }
	// rtsp.SnapshotWriter <- ivfFile
	render.Status(r, http.StatusOK)
}

//Record  录屏
func Record(w http.ResponseWriter, r *http.Request) {
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

//Ptz  控制
func Ptz(w http.ResponseWriter, r *http.Request) {
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
