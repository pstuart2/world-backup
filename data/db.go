package data

import (
	"time"

	"encoding/json"

	"os"

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
	fs   IDbFileSystem
	name string
	data dbData
}

type IDbFileSystem interface {
	Exists(path string) (bool, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

type IDb interface {
	Save() error
	Close()

	AddFolder(path string) *Folder
	Folders() []*Folder
	GetFolderByPath(path string) *Folder
}

func Open(name string, af IDbFileSystem) (*Db, error) {
	d, err := getData(name, af)
	if err != nil {
		return nil, err
	}

	return &Db{
		fs:   af,
		name: name,
		data: *d,
	}, nil
}

func getData(name string, af IDbFileSystem) (*dbData, error) {
	exists, err := af.Exists(name)
	if err != nil {
		return nil, err
	}

	if exists {
		return getDataFromFile(name, af)
	}

	d := dbData{
		CreatedAt: getNow(),
	}

	return &d, nil
}

func getDataFromFile(name string, af IDbFileSystem) (*dbData, error) {
	file, e := af.ReadFile(name)
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
