package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PeachIceTea/fela/conf"
	"github.com/PeachIceTea/fela/models"
)

func NewBook(r *gin.RouterGroup, c *conf.Config) {
	r.POST("/book/new", func(ctx *gin.Context) {
		b := models.Book{}
		err := ctx.BindJSON(&b)
		if err != nil {
			badRequest(ctx, "non json body")
			return
		}

		err = b.Insert(c)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, gin.H{"book_id": b.ID})
	})
}

func GetBooks(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/book", func(ctx *gin.Context) {
		books, err := models.GetBooks(c)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, gin.H{"books": books})
	})
}

func GetBook(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/book/:id", func(ctx *gin.Context) {
		id, err := idFromQuery(ctx)
		if err != nil {
			badRequest(ctx, err.Error())
			return
		}

		b, err := models.GetBook(id, c)
		if err != nil {
			if err == sql.ErrNoRows {
				notFound(ctx, fmt.Sprintf("no book with id %d", id))
				return
			}

			panic(err)
		}

		a, err := b.Audiobooks(c)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, gin.H{"book": b, "audiobooks": a})
	})
}
