package data

import (
	"time"
)

type Folder struct {
	Path       string    `json:"path"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

func (db *Db) AddFolder(path string) *Folder {
	f := Folder{
		Path:       path,
		CreatedAt:  getNow(),
		ModifiedAt: getNow(),
	}

	db.data.Folders = append(db.data.Folders, f)

	return &f
}
