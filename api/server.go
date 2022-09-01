package api

import (
	"dnslog/api/router"
	"dnslog/internal/response"
	"dnslog/web"
	"html/template"
	"io/fs"
	"net"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Addr   string
	engine *gin.Engine
	groups []router.Group
	l      net.Listener
}

func New(addr string) *Server {
	return &Server{
		Addr: addr,
	}
}

func (s *Server) accept() error {
	return nil
}

func (s *Server) Serve() error {
	addr, err := net.ResolveTCPAddr("tcp", s.Addr)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	gin.SetMode(gin.ReleaseMode)
	s.engine = gin.New()
	s.engine.Use(gin.Recovery())
	s.engine.Use(sessions.Sessions("dnslog", cookie.NewStore([]byte("secret"))))
	if gin.Mode() == gin.ReleaseMode {
		s.templateFS()
	}
	s.engine.SetTrustedProxies(nil)
	s.registerRoutes()
	go func() {
		s.engine.RunListener(listener)
	}()
	return nil
}

func (s *Server) templateFS() {
	tmpl := template.New("")
	staticFS, _ := fs.Sub(web.Static, "dist")
	tmpl = template.Must(tmpl.ParseFS(staticFS, "*.html"))
	s.engine.SetHTMLTemplate(tmpl)
	assetsFp, _ := fs.Sub(staticFS, "assets")
	s.engine.StaticFS("/assets", http.FS(assetsFp))
}

func (s *Server) index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func (s *Server) Group(relativePath string, routers []router.Router, handlers ...gin.HandlerFunc) {
	s.groups = append(s.groups, router.Group{
		RelativePath: relativePath,
		Routers:      routers,
		Handlers:     handlers,
	})
}

func (s *Server) registerRoutes() {
	for _, group := range s.groups {
		r := s.engine.Group(group.RelativePath, group.Handlers...)
		{
			for _, router := range group.Routers {
				s.register(r, router)
			}
		}
	}
	s.engine.GET("", s.index)
	s.engine.GET("/:path", s.index)
}

func (s *Server) register(r *gin.RouterGroup, router router.Router) {
	rg := r.Group(router.Prefix(), router.Middlewares()...)
	{
		for _, route := range router.Routes() {
			handler := s.makeHTTPHandler(route)
			rg.Handle(route.Method, route.Path, append(route.Middlewares, handler)...)
		}
	}
}

func (s *Server) makeHTTPHandler(route router.Route) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := (route.Action)(ctx)
		sessions.Default(ctx).Save()
		if err != nil {
			ctx.JSON(200, response.Error(err, data))
			return
		}
		ctx.JSON(200, response.Success(data))
	}
}
