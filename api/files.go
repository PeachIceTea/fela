package api

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/PeachIceTea/fela/conf"
)

// ServeFiles - /file/*path
func ServeFiles(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/files/*path", func(ctx *gin.Context) {
		auth := ctx.Query("auth")
		_, err := parseToken(auth, c)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				conf.M{"err": "could not parse auth token"},
			)
			return
		}

		path, _ := filepath.Abs(path.Clean(fmt.Sprintf("%s/%s", c.FilesPath, ctx.Param("path"))))
		http.ServeFile(ctx.Writer, ctx.Request, path)
	})
}
