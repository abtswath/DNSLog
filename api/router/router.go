package router

import "github.com/gin-gonic/gin"

type Router interface {
	Routes() []Route
	Prefix() string
	Middlewares() []gin.HandlerFunc
}
