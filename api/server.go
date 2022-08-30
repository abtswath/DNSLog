package api

import (
	"dnslog/api/router"
	"dnslog/internal/response"
	"dnslog/internal/store"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Listener net.Listener
	engine   *gin.Engine
	routers  []router.Router
	Store    store.Store
}

func (s *Server) Serve() error {
	var chError = make(chan error, 1)
	s.engine = gin.New()
	s.registerRoutes()
	go func() {
		logrus.Infof("HTTP server listen on: %s", s.Listener.Addr())
		chError <- s.engine.RunListener(s.Listener)
	}()
	return <-chError
}

func (s *Server) registerRoutes() {
	for _, router := range s.routers {
		s.register(router)
	}
}

func (s *Server) register(router router.Router) {
	s.engine.Group(router.Prefix(), router.Middlewares()...)
	{
		for _, route := range router.Routes() {
			handler := s.makeHTTPHandler(route)
			s.engine.Handle(route.Method, route.Path, append(route.Middlewares, handler)...)
		}
	}
}

func (s *Server) makeHTTPHandler(route router.Route) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := (route.Action)(ctx)
		if err != nil {
			ctx.JSON(200, response.Error(err, data))
			return
		}
		ctx.JSON(200, response.Success(data))
	}
}
