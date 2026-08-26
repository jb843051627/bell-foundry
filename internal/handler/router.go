package handler

import (
	"net/http"
	"strings"

	"github.com/jb843051627/bell-foundry/internal/service"
)

// Router 是轻量 API 路由器，保持业务服务与 HTTP 参数解析分离。
type Router struct {
	lab *service.Lab
}

// New 返回带恢复、请求 ID、CORS 和访问日志中间件的 HTTP 处理器。
func New(lab *service.Lab) http.Handler {
	router := &Router{lab: lab}
	return withCORS(withRequestID(withRecover(withAccessLog(router))))
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	if path == "" {
		path = "/"
	}
	switch {
	case path == "/":
		rt.handleHome(w, r)
	case strings.HasPrefix(path, "/assets/"):
		rt.handleHome(w, r)
	case strings.HasPrefix(path, "/assets/"):
		rt.handleHome(w, r)
	case path == "/healthz":
		rt.handleHealth(w, r)
	case strings.HasPrefix(path, "/api/specs"):
		rt.handleSpecs(w, r, path)
	case strings.HasPrefix(path, "/api/batches"):
		rt.handleBatches(w, r, path)
	case strings.HasPrefix(path, "/api/molds"):
		rt.handleMolds(w, r, path)
	case strings.HasPrefix(path, "/api/heats"):
		rt.handleHeats(w, r, path)
	case strings.HasPrefix(path, "/api/pours"):
		rt.handlePours(w, r, path)
	case strings.HasPrefix(path, "/api/curves"):
		rt.handleCurves(w, r, path)
	case strings.HasPrefix(path, "/api/bells"):
		rt.handleBells(w, r, path)
	case strings.HasPrefix(path, "/api/defects") || strings.HasPrefix(path, "/api/alerts"):
		rt.handleQuality(w, r, path)
	case path == "/api/report/daily":
		rt.handleDailyReport(w, r)
	default:
		http.NotFound(w, r)
	}
}
