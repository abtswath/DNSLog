package login

import (
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type loginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *loginRouter) login(ctx *gin.Context) (any, error) {
	var form loginForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if form.Username == l.username && form.Password == l.password {
		session := sessions.Default(ctx)
		session.Set("user", form.Username)
		return nil, nil
	}
	return nil, errors.New("用户名或密码错误")
}

func (l *loginRouter) logout(ctx *gin.Context) (any, error) {
	session := sessions.Default(ctx)
	session.Clear()
	return nil, nil
}
