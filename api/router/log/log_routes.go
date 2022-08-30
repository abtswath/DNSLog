package log

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (l *logRouter) getAll(ctx *gin.Context) (any, error) {
	offsetStr := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 30
	}
	return l.store.Paginate(uint(offset), uint(limit), ctx.Query("domain"), ctx.Query("ip")), nil
}
