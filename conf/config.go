package conf

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"text/template"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/hashicorp/go-envparse"
	"github.com/jmoiron/sqlx"
)

// M is a shortcut for the most often used map type.
type M map[string]interface{}

// Config stores app configuration. Expected to created once and handed around.
type Config struct {
	DB            *sqlx.DB
	ConnectionURL string
	Secret        []byte
	FilesPath     string
	Templates     *template.Template
}

// Initialize creates a new Config and runs all neccessary functions to get it
// up and running. This is the expected way to get a Config.
func Initialize() (c *Config) {
	c = &Config{}

	c.LoadEnv()
	c.EnsureDirectoryStructure()

	c.LoadTemplates()
	c.ConnectToDB()

	return
}

// EnsureDirectoryStructure ensures that the needed directory structure exists.
func (c *Config) EnsureDirectoryStructure() {
	c.ensureDirectory("")
	c.ensureDirectory("audio")
	c.ensureDirectory("cover")
}

func (c *Config) ensureDirectory(p string) error {
	err := os.Mkdir(path.Clean(fmt.Sprintf("%s/%s", c.FilesPath, p)), os.ModeDir|0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	return nil
}

// ConnectToDB connects to database using the ConnectionURL.
func (c *Config) ConnectToDB() {
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

// LoadEnv loads the .env file that should contain most of the configuration of
// a perticular fela instance.
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

	s, ok := env["SECRET"]
	if !ok {
		log.Panic("SECRET missing from .env")
	}
	c.Secret = []byte(s)

	if c.FilesPath, ok = env["FILES_PATH"]; !ok {
		log.Panic("FILES_PATH missing from .env")
	}
}

// LoadTemplates loads SQL templates.
func (c *Config) LoadTemplates() {
	c.Templates = template.Must(template.ParseGlob("db/templates/*.sql"))
}

// TemplateString returns a given template. Panics if template is not found.
func (c *Config) TemplateString(name string) string {
	buf := bytes.Buffer{}
	err := c.Templates.ExecuteTemplate(&buf, fmt.Sprintf("%s.sql", name), nil)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

// TemplateWithData returns a given template and accepts data to be passed to
// template. Panics if template is not found.
func (c *Config) TemplateWithData(name string, data interface{}) string {
	buf := bytes.Buffer{}
	err := c.Templates.ExecuteTemplate(&buf, fmt.Sprintf("%s.sql", name), &data)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}
