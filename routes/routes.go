package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Define errors
var (
	ErrIDParamMissing = errors.New("id parameter is missing")
	ErrInvalidID      = errors.New("invalid id")
)

// badRequest - Sends a bad request error
func badRequest(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
}

// notFound - Sends a not found error
func notFound(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusNotFound, gin.H{"error": msg})
}

// Extracts :id from URL
func idFromQuery(ctx *gin.Context) (int64, error) {
	idStr, ok := ctx.Params.Get("id")
	if !ok {
		return 0, ErrIDParamMissing
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, ErrInvalidID
	}

	return id, err
}
