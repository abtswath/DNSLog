package logs

import (
	"dnslog/internal/model"
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
	domain := ctx.Query("domain")
	ip := ctx.Query("ip")
	size, err := l.store.Size(domain, ip)
	if err != nil {
		return nil, err
	}
	var records []model.Record
	if size > 0 {
		records, err = l.store.Paginate(offset, limit, domain, ip)
		if err != nil {
			return nil, err
		}
	}
	return map[string]any{
		"data": records,
		"size": size,
	}, nil
}

func (l *logRouter) getByDomain(ctx *gin.Context) (any, error) {
	return l.store.GetByDomain(ctx.Param("domain"))
}
