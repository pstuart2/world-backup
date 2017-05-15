package api

import (
	"awana-app/server/data"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

const (
	tokenKey  = "app.token"
	loggerKey = "app.logger"
	dbKey     = "app.db"
)

func getLogger(ctx echo.Context) *logrus.Entry {
	obj := ctx.Get(loggerKey)
	if obj == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return obj.(*logrus.Entry)
}

func getToken(ctx echo.Context) *jwt.Token {
	obj := ctx.Get(tokenKey)
	if obj == nil {
		return nil
	}

	return obj.(*jwt.Token)
}

func getAppDb(ctx echo.Context) data.IAppDb {
	obj := ctx.Get(dbKey)
	if obj == nil {
		return nil
	}

	return obj.(data.IAppDb)
}
