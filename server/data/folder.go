package data

import (
	"path"
	"time"
)

type Folder struct {
	Id         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Path       string    `json:"path"`
	LastRun    time.Time `json:"lastRun"`
	Worlds     []*World  `json:"worlds"`
}

func (f *Folder) AddWorld(name string) *World {
	world := World{
		Id:        getId(),
		CreatedAt: getNow(),
		Name:      name,
		FullPath:  path.Join(f.Path, name),
	}

	f.Worlds = append(f.Worlds, &world)

	return &world
}

func (f *Folder) GetWorld(id string) *World {
	for i := range f.Worlds {
		if f.Worlds[i].Id == id {
			return f.Worlds[i]
		}
	}

	return nil
}

func (f *Folder) GetWorldByName(name string) *World {
	for i := range f.Worlds {
		if f.Worlds[i].Name == name {
			return f.Worlds[i]
		}
	}

	return nil
}

func (f *Folder) RemoveWorld(id string)  {
	for i := range f.Worlds {
		if f.Worlds[i].Id == id {
			copy(f.Worlds[i:], f.Worlds[i+1:])
			f.Worlds[len(f.Worlds)-1] = nil
			f.Worlds = f.Worlds[:len(f.Worlds)-1]
			break
		}
	}
}

func (db *Db) AddFolder(path string) *Folder {
	now := getNow()

	f := Folder{
		Id:         getId(),
		Path:       path,
		CreatedAt:  now,
		ModifiedAt: now,
	}

	db.data.Folders = append(db.data.Folders, &f)

	return &f
}

func (db *Db) Folders() []*Folder {
	return db.data.Folders
}

func (db *Db) GetFolderByPath(path string) *Folder {
	for i := range db.data.Folders {
		if db.data.Folders[i].Path == path {
			return db.data.Folders[i]
		}
	}

	return nil
}

func (db *Db) GetFolder(id string) *Folder {
	for i := range db.data.Folders {
		if db.data.Folders[i].Id == id {
			return db.data.Folders[i]
		}
	}

	return nil
}
