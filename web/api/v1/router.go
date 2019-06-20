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

	r.Get("/GetRunInfo", GetRunInfo)
	r.Get("/GetSystemInfo", GetSystemInfo)
	r.Get("/GetSrc", GetSrcInfo)
	r.Get("/Snapshot", Snapshot)
	r.Get("/Record", Record)
	r.Get("/Ptz", Ptz)
	return r
}
