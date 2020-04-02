package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/PeachIceTea/fela/conf"
)

// AuthorList - /author - Returns a list of all authors
func AuthorList(r *httprouter.Router, c *conf.Config) {
	r.GET("/author", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var authors []struct {
			Name string `db:"author"`
		}

		err := c.DB.Select(&authors, c.TemplateString("author_list"))
		if err != nil {
			fmt.Println(err)
			conf.JSONResponse(w, http.StatusInternalServerError, conf.M{"error": "internal server error"})
			return
		}

		authorList := make([]string, len(authors))
		for i, e := range authors {
			authorList[i] = e.Name
		}

		conf.JSONResponse(w, http.StatusOK, authorList)
	})
}
