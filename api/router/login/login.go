package login

import (
	"dnslog/api/router"

	"github.com/gin-gonic/gin"
)

type loginRouter struct {
	username string
	password string
	routes   []router.Route
}

func New(username, password string) router.Router {
	l := &loginRouter{
		username: username,
		password: password,
	}
	l.initRouter()
	return l
}

func (l *loginRouter) initRouter() {
	l.routes = []router.Route{
		router.NewPostRoute("/session", l.login),
		router.NewDeleteRoute("/session", l.logout),
	}
}

func (l *loginRouter) Routes() []router.Route {
	return l.routes
}

func (l *loginRouter) Prefix() string {
	return ""
}

func (l *loginRouter) Middlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}
