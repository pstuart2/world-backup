package data

import (
	"time"
)

type Folder struct {
	Id         string    `json:"id"`
	Path       string    `json:"path"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

func (db *Db) AddFolder(path string) *Folder {
	f := Folder{
		Id:         getId(),
		Path:       path,
		CreatedAt:  getNow(),
		ModifiedAt: getNow(),
	}

	db.data.Folders = append(db.data.Folders, f)

	return &f
}
