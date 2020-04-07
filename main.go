package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/PeachIceTea/fela/conf"
	"github.com/PeachIceTea/fela/routes"
)

func main() {
	c := conf.Init()
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:1234"}
	r.Use(cors.New(corsConfig))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "Hello, World!"})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
	})

	v1 := r.Group("/api/v1")
	{
		routes.GetBooks(v1, c)
		routes.GetBook(v1, c)
		routes.NewBook(v1, c)
		routes.NewAudiobook(v1, c)

		routes.Upload(v1, c)
		routes.AssignFiles(v1, c)
	}

	r.Run()
}
