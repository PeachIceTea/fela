package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/PeachIceTea/fela/conf"
)

// Book - /book/:id - Returns data for a single book
func Book(r *httprouter.Router, c *conf.Config) {
	r.GET("/book/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := strconv.Atoi(ps.ByName("id"))
		if err != nil || id <= 0 {
			conf.JSONResponse(w, http.StatusBadRequest, conf.M{"error": "invalid id"})
			return
		}

		var b conf.Book
		err = c.DB.Get(&b, c.TemplateString("book_select_by_id"), id)
		if err == sql.ErrNoRows {
			conf.JSONResponse(w, http.StatusNotFound, conf.M{"error": "book not found"})
			return
		} else if err != nil {
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			return
		}

		var f []conf.File
		err = c.DB.Select(&f, c.TemplateString("file_select_by_book"), b.ID)
		if err != nil {
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			return
		}

		conf.JSONResponse(w, http.StatusOK, conf.M{"book": b, "files": f})
	})
}
