package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/PeachIceTea/fela/conf"
	"github.com/PeachIceTea/fela/routes"
)

func main() {
	c := conf.Init()
	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		conf.JSONResponse(w, http.StatusOK, conf.M{"msg": "Hello, World!"})
	})

	routes.Upload(r, &c)
	routes.BookCreate(r, &c)
	routes.BookList(r, &c)
	routes.Book(r, &c)

	r.ServeFiles("/files/*filepath", http.Dir(c.FilePath))

	n := negroni.Classic()
	n.Use(cors.New(cors.Options{}))
	n.UseHandler(r)

	n.Run(":8080")
}
