package api

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"fmt"
	"world-backup/server/conf"
	"world-backup/server/data"
)

var getNow = time.Now

type IApiDb interface {
	Folders() []*data.Folder
	GetFolder(id string) *data.Folder
	Save() error
}

type IApiFileSystem interface {
	Exists(path string) (bool, error)
	Remove(name string) error
	Unzip(src, dest string) error
	Rename(oldname, newname string) error
}

// API is the data holder for the API
type API struct {
	log    *logrus.Entry
	config *conf.Config
	Server IServer
	Db     IApiDb
	Fs     IApiFileSystem
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// Start will start the API on the specified port
func (api *API) Start() error {
	return api.Server.Start(fmt.Sprintf(":%d", api.config.Port))
}

// NewAPI will create an api instance that is ready to start
func NewAPI(log *logrus.Entry, config *conf.Config, db IApiDb, fs IApiFileSystem) *API {
	echoServer := EchoServer{e: echo.New()}

	// create the api
	api := &API{
		config: config,
		log:    log.WithField("component", "api"),
		Server: echoServer,
		Db:     db,
		Fs:     fs,
	}

	return api
}
