package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/PeachIceTea/fela/conf"
)

// BookCreate - POST book/create - Creates a new book entry and links a file
func BookCreate(r *httprouter.Router, c *conf.Config) {
	r.POST("/book/create", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var data struct {
			Title       string `json:"title"`
			Author      string `json:"author"`
			Description string `json:"description"`
		}

		err := conf.JSONBody(&data, r)
		if err != nil {
			conf.JSONResponse(w, http.StatusBadRequest, conf.M{"error": "not valid json"})
			return
		}

		if data.Title == "" || data.Author == "" {
			conf.JSONResponse(w, http.StatusBadRequest, conf.M{"error": "data missing"})
			return
		}

		bookID, err := (func() (bookID int64, err error) {
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

			bookID, err = res.LastInsertId()
			if err != nil {
				tx.Rollback()
				return
			}

			return
		})()
		if err != nil {
			fmt.Println(err)
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			return
		}

		conf.JSONResponse(w, http.StatusOK, conf.M{"book_id": bookID})
	})
}
