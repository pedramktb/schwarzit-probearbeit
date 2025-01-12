package ginRouter

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/pedramktb/schwarzit-probearbeit/internal/dtos"
	"github.com/pedramktb/schwarzit-probearbeit/internal/logging"
	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
)

func ErrorResponse(c *gin.Context, err error) {
	switch {
	case errors.IsAny(err, types.ErrNotFound):
		logging.FromContext(c.Request.Context()).Debug("Not found", zap.Error(err))
		c.JSON(http.StatusNotFound, dtos.ErrorResponse{Error: err.Error()})
	case errors.IsAny(err, types.ErrBadRequest):
		logging.FromContext(c.Request.Context()).Debug("Bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Error: err.Error()})
	case errors.IsAny(err, types.ErrUnauthorized):
		logging.FromContext(c.Request.Context()).Debug("Unauthorized", zap.Error(err))
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Error: err.Error()})
	case errors.IsAny(err, types.ErrForbidden):
		logging.FromContext(c.Request.Context()).Debug("Forbidden", zap.Error(err))
		c.JSON(http.StatusForbidden, dtos.ErrorResponse{Error: err.Error()})
	case errors.IsAny(err, types.ErrInternal):
		logging.FromContext(c.Request.Context()).Warn("Internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: err.Error()})
	default:
		logging.FromContext(c.Request.Context()).Error("Unknown error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: err.Error()})
	}
}
