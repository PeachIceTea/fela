package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/PeachIceTea/fela/conf"
	"github.com/PeachIceTea/fela/models"
	"github.com/gin-gonic/gin"
)

func NewAudiobook(r *gin.RouterGroup, c *conf.Config) {
	r.POST("/audiobook/new", func(ctx *gin.Context) {
		var data struct {
			BookID int64 `json:"book_id"`
		}
		err := ctx.BindJSON(&data)
		if err != nil {
			badRequest(ctx, "non json body")
			return
		}

		a := models.Audiobook{
			Book: data.BookID,
		}
		err = a.Insert(c)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, gin.H{"audiobook_id": a.ID})
	})
}

func GetAudiobook(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/audiobook/:id", func(ctx *gin.Context) {
		id, err := idFromQuery(ctx)
		if err != nil {
			badRequest(ctx, err.Error())
			return
		}

		a, err := models.GetAudiobook(id, c)
		if err != nil {
			if err == sql.ErrNoRows {
				notFound(ctx, fmt.Sprintf("no book with id %d", id))
				return
			}

			panic(err)
		}

		fl, err := a.Files(c)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, gin.H{"audiobooks": a, "files": fl})
	})
}
