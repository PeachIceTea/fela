package models

import "github.com/PeachIceTea/fela/conf"

type Audiobook struct {
	ID int64 `db:"id" json:"id"`

	Book int64 `db:"book" json:"book"`

	CreatedAt string `db:"created_at" json:"created_at"`
}

func (a *Audiobook) Insert(c *conf.Config) (err error) {
	res, err := c.DB.Exec(c.TemplateString("audiobook_insert"), a.Book)
	if err != nil {
		return
	}

	a.ID, err = res.LastInsertId()
	return
}

type AudiobookAssignment struct {
	FileID  int64 `json:"file_id"`
	Chapter int64 `json:"chapter"`
}

func (a *Audiobook) AssignFiles(assignments []AudiobookAssignment, c *conf.Config) (err error) {
	//TODO: Check if Audiobook and Files exist

	tx, err := c.DB.Beginx()
	if err != nil {
		return
	}
	defer tx.Commit()

	sqlStr := c.TemplateString("file_assign_audiobook")
	for _, as := range assignments {
		_, err = tx.Exec(sqlStr, a.ID, as.Chapter, as.FileID)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	return
}
