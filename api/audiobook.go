package api

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/PeachIceTea/fela/conf"
	"github.com/gin-gonic/gin"
)

// Audiobook represents a Audiobook database row.
type Audiobook struct {
	ID        int64   `db:"id" json:"id"`
	Title     *string `db:"title" json:"title"`
	Author    *string `db:"author" json:"author"`
	Uploader  *int64  `db:"uploader" json:"uploader"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt *string `db:"updated_at" json:"updated_at"`

	Files []File `db:"-" json:"files,omitempty"`
}

// GetAudiobooks - GET /audiobook - Get all audiobooks.
func GetAudiobooks(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/audiobook", func(ctx *gin.Context) {
		audiobooks := []Audiobook{}
		err := c.DB.Select(&audiobooks, c.TemplateString("all_audiobooks"))
		if err != nil {
			panic(err)
		}
		ctx.JSON(http.StatusOK, conf.M{"audiobooks": audiobooks})
	})
}

// GetAudiobook - GET /audiobook/:id - Get single audiobook including its files.
func GetAudiobook(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/audiobook/:id", func(ctx *gin.Context) {
		id, err := getID(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, conf.M{"err": err.Error()})
			return
		}

		audiobook := Audiobook{}
		err = c.DB.Get(&audiobook, c.TemplateString("get_audiobook"), id)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(
					http.StatusNotFound,
					conf.M{"err": "no audiobook with that id"},
				)
				return
			}

			panic(err)
		}

		err = c.DB.Select(
			&audiobook.Files, c.TemplateString("get_audiobook_files"),
			audiobook.ID)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, conf.M{"audiobook": audiobook})
	})
}

// GetAudiobookFiles - GET /audiobook/:id/files - Get files for audiobook.
func GetAudiobookFiles(r *gin.RouterGroup, c *conf.Config) {
	r.GET("/audiobook/:id/files", func(ctx *gin.Context) {
		id, err := getID(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, conf.M{"err": err.Error()})
			return
		}

		files := []File{}
		err = c.DB.Select(&files, c.TemplateString("get_audiobook_files"), id)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, conf.M{"files": files})
	})
}

// UpdateAudiobook - PUT /audiobook/:id - Updates audiobook.
// Accepts the fields "title" and "author" and "cover" as file upload.
func UpdateAudiobook(r *gin.RouterGroup, c *conf.Config) {
	r.PUT("/audiobook/:id", func(ctx *gin.Context) {
		id, err := getID(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, conf.M{"err": err.Error()})
			return
		}

		claims := getClaims(ctx)
		if !claims.isUploader() {
			ctx.JSON(
				http.StatusUnauthorized,
				conf.M{"err": "no permission to update audiobook"},
			)
			return
		}

		var data struct {
			ID     int64                    `form:"-" json:"-" db:"id"`
			Title  *string                  `form:"title" json:"title"`
			Author *string                  `form:"author" json:"author"`
			Cover  *[]*multipart.FileHeader `form:"cover"`
		}

		err = ctx.ShouldBind(&data)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				conf.M{"err": "invalid request body"},
			)
			return
		}

		if data.Title == nil && data.Author == nil && data.Cover == nil {
			ctx.JSON(
				http.StatusBadRequest,
				conf.M{"err": "no fields to update"},
			)
			return
		}

		if data.Cover != nil {
			if len(*data.Cover) != 1 {
				ctx.JSON(
					http.StatusBadRequest,
					conf.M{"err": "only one cover upload allowed"},
				)
				return
			}

			f, err := (*data.Cover)[0].Open()
			if err != nil {
				panic(err)
			}
			defer f.Close()

			disk, err := os.Create(
				path.Clean(fmt.Sprintf("%s/cover/%d.jpg", c.FilesPath, id)),
			)
			if err != nil {
				panic(err)
			}

			_, err = io.Copy(disk, f)
			if err != nil {
				panic(err)
			}

			err = disk.Close()
			if err != nil {
				panic(err)
			}
		}

		if data.Title != nil || data.Author != nil {
			data.ID = id
			_, err = c.DB.NamedExec(
				c.TemplateWithData("update_audiobook", data),
				data,
			)
			if err != nil {
				panic(err)
			}
		}

		ctx.JSON(http.StatusOK, conf.M{"msg": "audiobook updated"})
	})
}

// DeleteAudiobook - DELETE /audiobook/:id - Deletes audiobook and files.
func DeleteAudiobook(r *gin.RouterGroup, c *conf.Config) {
	r.DELETE("/audiobook/:id", func(ctx *gin.Context) {
		id, err := getID(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, conf.M{"err": err.Error()})
			return
		}

		claims := getClaims(ctx)
		if !claims.isUploader() {
			ctx.JSON(
				http.StatusUnauthorized,
				conf.M{"err": "no permission to delete"},
			)
			return
		}

		_, err = c.DB.Exec(c.TemplateString("delete_audiobook_files"), id)
		if err != nil {
			panic(err)
		}
		_, err = c.DB.Exec(c.TemplateString("delete_audiobook"), id)
		if err != nil {
			panic(err)
		}

		// Remove audio files and cover
		os.RemoveAll(path.Clean(fmt.Sprintf("%s/audio/%d", c.FilesPath, id)))
		os.Remove(path.Clean(fmt.Sprintf("%s/cover/%d.jpg", c.FilesPath, id)))

		ctx.JSON(http.StatusOK, conf.M{"msg": "audiobook deleted"})
	})
}
