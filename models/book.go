package models

import (
	"github.com/PeachIceTea/fela/conf"
)

type Book struct {
	ID int64 `db:"id" json:"id"`

	Title  string `db:"title" json:"title"`
	Author string `db:"author" json:"author"`

	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt *string `db:"updated_at" json:"updated_at"`
}

func (b *Book) Insert(c *conf.Config) (err error) {
	//TODO: Rewrite as normal function
	res, err := c.DB.Exec(c.TemplateString("book_insert"), b.Title, b.Author)
	if err != nil {
		return
	}

	b.ID, err = res.LastInsertId()
	return
}

func (b *Book) Audiobooks(c *conf.Config) (a []*Audiobook, err error) {
	err = c.DB.Select(&a, c.TemplateString("audiobook_select_by_book"), b.ID)
	return
}

func GetBooks(c *conf.Config) (b []*Book, err error) {
	err = c.DB.Select(&b, c.TemplateString("book_select_all"))
	return
}

func GetBook(id int64, c *conf.Config) (b *Book, err error) {
	b = &Book{}
	err = c.DB.Get(b, c.TemplateString("book_select_by_id"), id)
	return
}
