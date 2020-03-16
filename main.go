package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func main() {
	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		JSONResponse(w, 200, m{"msg": "Hello, World!"})
	})

	n := negroni.Classic()
	n.Use(cors.New(cors.Options{}))
	n.UseHandler(r)

	n.Run(":8080")
}

type m map[string]interface{}

// JSONBody - Get JSON Body from request
func JSONBody(e interface{}, r *http.Request) (err error) {
	err = json.NewDecoder(r.Body).Decode(e)
	return
}

// JSONResponse - Respond with JSON
func JSONResponse(w http.ResponseWriter, status int, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(b)
	if err != nil {
		panic(err)
	}
}
