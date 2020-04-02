package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/PeachIceTea/fela/conf"
)

// BookList - /book - Returns a list of all books
func BookList(r *httprouter.Router, c *conf.Config) {
	r.GET("/book", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var book []conf.Book

		err := c.DB.Select(&book, c.TemplateString("book_all"))
		if err != nil {
			fmt.Println(err)
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			return
		}

		conf.JSONResponse(w, http.StatusOK, &book)
	})
}
