package data

import (
	"time"

	"encoding/json"

	"github.com/spf13/afero"
	"github.com/ventu-io/go-shortid"
)

var getNow = time.Now
var getId = shortid.MustGenerate

type dbData struct {
	CreatedAt time.Time `json:"createdAt"`
	LastSave  time.Time `json:"lastSave"`
	Folders   []*Folder `json:"folders"`
}

type Db struct {
	fs   afero.Afero
	name string
	data dbData
}

type IDb interface {
	Save() error
	Close()

	AddFolder(path string) *Folder
	Folders() []*Folder
	GetFolderByPath(path string) *Folder
}

func Open(name string, fs afero.Fs) *Db {
	af := afero.Afero{Fs: fs}
	d, _ := getData(name, af)

	return &Db{
		fs:   af,
		name: name,
		data: *d,
	}
}

func getData(name string, fs afero.Afero) (*dbData, error) {
	exists, err := fs.Exists(name)
	if err != nil {
		return nil, err
	}

	if exists {
		return getDataFromFile(name, fs)
	}

	d := dbData{
		CreatedAt: getNow(),
	}

	return &d, nil
}

func getDataFromFile(name string, fs afero.Afero) (*dbData, error) {
	file, e := fs.ReadFile(name)
	if e != nil {
		return nil, e
	}

	var jsonData dbData
	if err := json.Unmarshal(file, &jsonData); err != nil {
		return nil, err
	}

	return &jsonData, nil
}

func (db *Db) Save() error {
	db.data.LastSave = getNow()

	jsonData, err := json.Marshal(db.data)
	if err != nil {
		return err
	}

	if wErr := db.fs.WriteFile(db.name, jsonData, 0600); wErr != nil {
		return wErr
	}

	return nil
}

func (db *Db) Close() {
	db.Save()
}
