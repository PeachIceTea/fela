// The main package creates the gin server and a conf.Config. It hands off
// creating routes to the api package. It also contains a slightly modified
// version of the default gin panic middleware.

package main

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/PeachIceTea/fela/api"
	"github.com/PeachIceTea/fela/conf"
)

func main() {
	c := conf.Initialize()
	r := gin.New()

	r.Use(recoverMiddleware())
	r.Use(static.Serve("/", static.LocalFile("./client/dist", false)))
	r.NoRoute(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/api") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
			return
		}

		http.ServeFile(ctx.Writer, ctx.Request, "./client/dist/index.html")
	})

	api.RegisterRoutes(r.Group("/"), c)

	r.Run()
}
