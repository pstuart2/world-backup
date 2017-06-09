package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"fmt"
	"world-backup/conf"
	"world-backup/data"
)

type IApiDb interface {
	Folders() []*data.Folder
}

// API is the data holder for the API
type API struct {
	log    *logrus.Entry
	config *conf.Config
	Server IServer
	Db     IApiDb
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// Start will start the API on the specified port
func (api *API) Start() error {
	return api.Server.Start(fmt.Sprintf(":%d", api.config.Port))
}

// NewAPI will create an api instance that is ready to start
func NewAPI(log *logrus.Entry, config *conf.Config, db IApiDb) *API {
	echoServer := EchoServer{e: echo.New()}

	// create the api
	api := &API{
		config: config,
		log:    log.WithField("component", "api"),
		Server: echoServer,
		Db:     db,
	}

	return api
}
