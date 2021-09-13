package api

import (
	"net/http"

	"github.com/PeachIceTea/fela/conf"
	"github.com/gin-gonic/gin"
)

// Progress represents a Progress database row.
type Progress struct {
	ID        int64   `db:"id" json:"id"`
	User      int64   `db:"user" json:"user"`
	Audiobook int64   `db:"audiobook" json:"audiobook"`
	File      int64   `db:"file" json:"file"`
	Progress  float64 `db:"progress" json:"progress"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt *string `db:"updated_at" json:"updated_at"`
}

// UpdateProgress
// TODO: Check if file is part of audiobook
func UpdateProgress(r *gin.RouterGroup, c *conf.Config) {
	r.PUT("/progress", func(ctx *gin.Context) {
		claims := getClaims(ctx)

		var data struct {
			User      int64   `form:"-" json:"-"`
			Audiobook int64   `form:"audiobook" json:"audiobook"`
			File      int64   `form:"file" json:"file"`
			Progress  float64 `form:"progress" json:"progress"`
		}
		err := ctx.ShouldBind(&data)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				conf.M{"err": "invalid request body"},
			)
			return
		}

		data.User = claims.ID
		_, err = c.DB.NamedExec(
			c.TemplateString("update_progress"),
			data,
		)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, conf.M{"msg": "progress updated"})
	})
}

// GetProgress returns the progress of all audiobooks a user has listened to.
func GetProgress(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/progress", func(ctx *gin.Context) {
		id := getClaims(ctx).ID

		progressList := []Progress{}
		err := c.DB.Select(&progressList, c.TemplateString("get_progress_user"), id)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, conf.M{"progress": progressList})
	})
}
