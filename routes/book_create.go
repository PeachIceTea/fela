package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/PeachIceTea/fela/conf"
)

// BookCreate - /book/create - Creates a new book entry and links a file
func BookCreate(r *httprouter.Router, c *conf.Config) {
	r.POST("/book/create", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var data struct {
			Title       string `json:"title"`
			Author      string `json:"author"`
			Description string `json:"description"`
			FileID      int64  `json:"file_id"`
		}

		err := conf.JSONBody(&data, r)
		if err != nil {
			conf.JSONResponse(w, http.StatusBadRequest, conf.M{"error": "not valid json"})
			return
		}

		if data.Title == "" || data.Author == "" || data.FileID == 0 {
			conf.JSONResponse(w, http.StatusBadRequest, conf.M{"error": "data missing"})
			return
		}

		unassigned, err := (func() (u bool, err error) {
			var file struct {
				Book *int64 `db:"book"`
			}

			err = c.DB.Get(&file, c.TemplateWithData("file_select_by_id", []string{"book"}), data.FileID)
			if err != nil {
				return
			}

			u = file.Book == nil
			return
		})()
		if err != nil {
			conf.JSONResponse(w, http.StatusBadRequest, conf.M{"error": "cannot find file"})
			return
		}
		if !unassigned {
			//TODO: Enable option to reassign books
			conf.JSONResponse(w, http.StatusBadRequest, conf.M{"error": "file already assigned"})
			return
		}

		err = (func() (err error) {
			tx, err := c.DB.Beginx()
			if err != nil {
				return
			}
			defer tx.Commit()

			res, err := tx.Exec(c.TemplateString("book_insert"), data.Author, data.Title, data.Description)
			if err != nil {
				tx.Rollback()
				return
			}

			bookID, err := res.LastInsertId()
			if err != nil {
				tx.Rollback()
				return
			}

			_, err = tx.Exec(c.TemplateString("file_assign_book"), bookID, data.FileID)
			if err != nil {
				return
			}

			return
		})()
		if err != nil {
			fmt.Println(err)
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			return
		}

		conf.JSONResponse(w, http.StatusOK, conf.M{"msg": "book has been created"})
	})
}
