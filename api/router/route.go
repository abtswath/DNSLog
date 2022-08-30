package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ActionFunc = func(ctx *gin.Context) (any, error)

type Route struct {
	Path        string
	Method      string
	Action      ActionFunc
	Middlewares []gin.HandlerFunc
}

func NewRoute(path, method string, action ActionFunc, middlewares ...gin.HandlerFunc) Route {
	return Route{
		Path:        path,
		Method:      method,
		Action:      action,
		Middlewares: middlewares,
	}
}

func NewGetRoute(path string, action ActionFunc, middlewares ...gin.HandlerFunc) Route {
	return NewRoute(path, http.MethodGet, action, middlewares...)
}
