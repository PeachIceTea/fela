package routes

import (
	"mime/multipart"
	"net/http"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/PeachIceTea/fela/conf"
	"github.com/PeachIceTea/fela/models"
)

// Upload - POST /file/upload - Upload a file
func Upload(r *gin.RouterGroup, c *conf.Config) {
	r.POST("/file/upload", func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			badRequest(ctx, "non multipart/form-data body")
			return
		}

		files, ok := form.File["file"]
		if !ok {
			badRequest(ctx, "file missing")
			return
		}

		if len(files) > 1 {
			sort.Slice(files, func(i, j int) bool {
				return files[i].Filename < files[j].Filename
			})
		}

		var wg sync.WaitGroup
		type fileResponse struct {
			Name string `json:"name"`
			ID   int64  `json:"id"`
		}
		ids := make([]fileResponse, len(files))
		for i, f := range files {
			wg.Add(1)
			go func(i int, f *multipart.FileHeader) {
				mf := models.File{Name: f.Filename}

				fs, err := f.Open()
				if err != nil {
					//TODO: Figure out the best thing to do when a file fails
					//TODO: Delete already saved files on error
					panic(err)
				}

				err = mf.Insert(fs, c)
				if err != nil {
					panic(err)
				}

				ids[i] = fileResponse{f.Filename, mf.ID}
				wg.Done()
			}(i, f)
		}
		wg.Wait()

		ctx.JSON(http.StatusOK, gin.H{"file_ids": ids})
	})
}

func AssignFiles(r *gin.RouterGroup, c *conf.Config) {
	r.POST("/file/assign", func(ctx *gin.Context) {
		var data struct {
			AudiobookID int64                        `json:"audiobook_id"`
			Assignments []models.AudiobookAssignment `json:"assignments"`
		}
		err := ctx.BindJSON(&data)
		if err != nil {
			badRequest(ctx, "non json body")
			return
		}

		a := models.Audiobook{ID: data.AudiobookID}
		err = a.AssignFiles(data.Assignments, c)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, gin.H{"msg": "files assigned"})
	})
}
