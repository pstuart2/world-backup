package api

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pborman/uuid"
)

func (api *API) SetUpRoutes() {

	api.Server.Use(api.setupRequest)
	api.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	api.Server.Use(middleware.Static(api.config.StaticRoot))
	api.Server.GET("*", api.index)

	apiGroup := api.Server.Group("/api")
	apiGroup.GET("/folders", api.getFolders)
	apiGroup.GET("/folders/:id/worlds", api.getWorlds)
	apiGroup.DELETE("/folders/:id/worlds/:wid/backups/:bid", api.deleteWorldBackup)
	apiGroup.PATCH("/folders/:id/worlds/:wid/backups/:bid", api.restoreWorldBackup)

	routes := api.Server.Routes()
	for i := 0; i < len(routes); i++ {
		api.log.Info(routes[i].Method + ": " + routes[i].Path)
	}
}

func (api *API) setupRequest(f echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		logger := api.log.WithFields(logrus.Fields{
			"method":     req.Method,
			"path":       req.URL.Path,
			"request_id": uuid.NewRandom().String(),
		})
		ctx.Set(loggerKey, logger)

		startTime := time.Now()
		defer func() {
			rsp := ctx.Response()
			logger.WithFields(logrus.Fields{
				"status_code":  rsp.Status,
				"runtime_nano": time.Since(startTime).Nanoseconds(),
			}).Info("Finished request")
		}()

		logger.WithFields(logrus.Fields{
			"user_agent":     req.UserAgent(),
			"content_length": req.ContentLength,
		}).Info("Starting request")

		// we have to do this b/c if not the final error handler will not
		// in the chain of middleware. It will be called after meaning that the
		// response won't be set properly.
		err := f(ctx)
		if err != nil {
			ctx.Error(err)
		}
		return err
	}
}

func (api *API) index(ctx echo.Context) error {
	log := getLogger(ctx)
	log.Infof("Returning index from: %s", api.config.StaticRoot)
	return ctx.File(api.config.StaticRoot + "index.html")
}
