package logs

import (
	"dnslog/api/router"
	"dnslog/internal/store"

	"github.com/gin-gonic/gin"
)

type logRouter struct {
	store  store.Store
	routes []router.Route
}

func New(store store.Store) router.Router {
	l := &logRouter{
		store: store,
	}
	l.initRouter()
	return l
}

func (l *logRouter) initRouter() {
	l.routes = []router.Route{
		router.NewGetRoute("/logs", l.getAll),
		router.NewGetRoute("/logs/:domain", l.getByDomain),
	}
}

func (l *logRouter) Routes() []router.Route {
	return l.routes
}

func (l *logRouter) Prefix() string {
	return ""
}

func (l *logRouter) Middlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}
