package models

import "github.com/PeachIceTea/fela/conf"

func GetAuthors(c *conf.Config) (a []string, err error) {
	var data []*Book
	err = c.DB.Select(&data, c.TemplateString("book_select_all_authors"))

	a = make([]string, len(data))
	for i, b := range data {
		a[i] = b.Author
	}
	return
}
