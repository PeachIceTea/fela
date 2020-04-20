// Fela is a audiobook management server
package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/PeachIceTea/fela/api"
	"github.com/PeachIceTea/fela/conf"
)

func main() {
	c := conf.Init()
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(recoverMiddleware())
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:1234"}
	corsConfig.AllowHeaders = []string{"content-type", "authorization"}
	r.Use(cors.New(corsConfig))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "Hello, World!"})
	})

	r.Use(static.Serve("/files", static.LocalFile("files", false)))

	api.RegisterRoutes(r.Group("/"), c)
	r.Run()
}
