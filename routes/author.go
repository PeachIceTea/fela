package routes

import (
	"net/http"

	"github.com/PeachIceTea/fela/conf"
	"github.com/PeachIceTea/fela/models"
	"github.com/gin-gonic/gin"
)

func GetAuthors(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/author", func(ctx *gin.Context) {
		authors, err := models.GetAuthors(c)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, gin.H{"authors": authors})
	})
}
