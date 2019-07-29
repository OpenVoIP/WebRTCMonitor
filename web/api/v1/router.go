package v1

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

//Router 路由
func Router() chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	// r.Use(AdminOnly)

	r.Get("/GetRunInfo", GetRunInfo)       // 运行信息
	r.Get("/GetSystemInfo", GetSystemInfo) //系统信息
	r.Get("/GetSrc", GetSrcInfo)           // 从配置文件读取源
	r.Get("/Snapshot", Snapshot)           //截图
	r.Get("/Record", Record)               // 回放
	r.Get("/Ptz", Ptz)                     // 控制指令
	return r
}
