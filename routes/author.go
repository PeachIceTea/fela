package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/PeachIceTea/fela/conf"
)

// Author - /author/:name - Returns a list of all books written by given author
func Author(r *httprouter.Router, c *conf.Config) {
	r.GET("/author/:name", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var books []conf.Book

		name := ps.ByName("name")

		err := c.DB.Select(&books, c.TemplateString("book_select_by_author_name"), name)
		if err != nil {
			fmt.Println(err)
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			return
		}

		conf.JSONResponse(w, http.StatusOK, books)
	})
}
