package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const authHeader = "X-Auth"

func Auth(tokens []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.GetHeader(authHeader)
		if authToken != "" {
			for _, token := range tokens {
				if token == authToken {
					ctx.Next()
					return
				}
			}
			ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Unauthorized."))
		} else {
			session := sessions.Default(ctx)
			defer session.Save()
			user := session.Get("user")
			if user == nil {
				if strings.Contains(ctx.GetHeader("Accept"), "application/json") {
					ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Unauthorized."))
				} else {
					ctx.Redirect(http.StatusFound, "/login")
				}
				return
			}
		}
		ctx.Next()
	}
}
