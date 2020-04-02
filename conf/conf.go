package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/hashicorp/go-envparse"
	"github.com/jmoiron/sqlx"
)

// Config - Stores app configuration
type Config struct {
	DB            *sqlx.DB
	ConnectionURL string
	Secret        string
	FilePath      string
	Templates     *template.Template
}

// Init - Initialize the server config
func Init() (c Config) {
	c.LoadEnv()
	c.LoadTemplates()
	c.ConnectToDatabase()

	return
}

// ConnectToDatabase - Connects to database using the ConnectionURL
func (c *Config) ConnectToDatabase() {
	u, err := url.Parse(c.ConnectionURL)
	if err != nil {
		log.Panic(err)
	}

	c.DB, err = sqlx.Open(
		u.Scheme,
		fmt.Sprintf("%s@%s%s", u.User, u.Host, u.Path),
	)
	if err != nil {
		log.Panic(err)
	}
}

// LoadEnv - Loads .env file
func (c *Config) LoadEnv() {
	f, err := os.Open(".env")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	env, err := envparse.Parse(f)
	if err != nil {
		log.Panic(err)
	}

	var ok bool
	if c.ConnectionURL, ok = env["DATABASE_URL"]; !ok {
		log.Panic("DATABASE_URL missing from .env")
	}

	if c.Secret, ok = env["SECRET"]; !ok {
		log.Panic("SECRET missing from .env")
	}

	if c.FilePath, ok = env["FILE_PATH"]; !ok {
		log.Panic("FILE_PATH missing from .env")
	}
}

// LoadTemplates - Loads SQL templates
func (c *Config) LoadTemplates() {
	c.Templates = template.Must(template.ParseGlob("db/templates/*.sql"))
}

// TemplateString - Returns after executing template
func (c *Config) TemplateString(name string) string {
	buf := bytes.Buffer{}
	err := c.Templates.ExecuteTemplate(&buf, fmt.Sprintf("%s.sql", name), nil)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (c *Config) TemplateWithData(name string, data interface{}) string {
	buf := bytes.Buffer{}
	err := c.Templates.ExecuteTemplate(&buf, fmt.Sprintf("%s.sql", name), &data)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

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
