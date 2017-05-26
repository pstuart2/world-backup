package data

import (
	"time"
)

type Folder struct {
	Id         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Path       string    `json:"path"`
	LastBackup time.Time `json:"lastBackup"`
}

func (db *Db) AddFolder(path string) *Folder {
	now := getNow()

	f := Folder{
		Id:         getId(),
		Path:       path,
		CreatedAt:  now,
		ModifiedAt: now,
	}

	db.data.Folders = append(db.data.Folders, f)

	return &f
}

func (db *Db) Folders() []Folder {
	return db.data.Folders
}

func (db *Db) FolderByPath(path string) *Folder {
	for i := range db.data.Folders {
		if db.data.Folders[i].Path == path {
			return &db.data.Folders[i]
		}
	}

	return nil
}