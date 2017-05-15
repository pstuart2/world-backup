package api

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/pborman/uuid"
)

func (api *API) SetUpRoutes() {

	// TODO: No server yet but will have one in the future.

	api.Server.Use(api.setupRequest)
	//api.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"localhost"},
	//}))

	//api.Server.Use(middleware.Static(api.config.StaticRoot))
	//api.Server.GET("*", api.index)

	//api.Server.POST("/auth/signup", api.signUp)
	//api.Server.POST("/auth/login", api.login)
	//
	//secureGroup := api.Server.Group("/api", requireClaims)
	//secureGroup.GET("/authenticate", api.authenticate)
	//secureGroup.GET("/clubTypes", api.clubTypes)
	//secureGroup.GET("/awards/:clubTypeId", api.awards)
	//
	//secureGroup.POST("/clubs", api.createClub)
	//secureGroup.GET("/clubs", api.clubs)
	//
	//secureGroup.POST("/books", api.createBook)
	//secureGroup.PATCH("/books/:id", api.updateBook)
	//secureGroup.GET("/books", api.books)
	//secureGroup.GET("/books/:id", api.book)
	//
	//secureGroup.POST("/books/:bookId/chapters", api.addChapter)
	//secureGroup.POST("/books/:bookId/chapters/:chapterId", api.updateChapter)
	//secureGroup.DELETE("/books/:bookId/chapters/:chapterId", api.deleteChapter)

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
	return nil // ctx.File(api.config.StaticRoot + "index.html")
}
