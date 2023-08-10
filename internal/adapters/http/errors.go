package http

import (
	"errors"
	"net/http"
	"url-shortener/internal/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
)

func (a *Adapter) ErrorHandler(ctx *gin.Context, err error) {
	l := zapctx.Logger(ctx)
	l.Sugar().Errorf("request failed: %s", err.Error())

	switch {
	case errors.Is(err, models.ErrNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
	case errors.Is(err, models.ErrStorage), errors.Is(err, models.ErrKeyGenerationFailed):
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
