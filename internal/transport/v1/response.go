package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/o-sokol-o/evo-fintech/internal/domain"
)

func newErrorResponse(ctx *gin.Context, statusCode int, err error) {
	logrus.Error(err)
	ctx.AbortWithStatusJSON(statusCode, domain.ErrorResponse{Error: err.Error()})
}
