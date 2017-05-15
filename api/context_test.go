package api

import (
	"testing"

	"net/http"
	"net/http/httptest"
	"strings"

	"awana-app/server/data"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestContext(t *testing.T) {

	Convey("Given an echo.Context", t, func() {
		e := echo.New()
		req, _ := http.NewRequest(echo.GET, "/some/path", strings.NewReader(""))
		rec := httptest.NewRecorder()

		Convey("When that context does not have a logger", func() {
			c := e.NewContext(req, rec)
			foundLogger := getLogger(c)

			Convey("getLogger should return the logger", func() {
				So(foundLogger, ShouldNotBeNil)
			})
		})

		Convey("When that context has a logger", func() {
			c := e.NewContext(req, rec)

			logger := logrus.WithField("test", "TestGetLogger1")
			c.Set(loggerKey, logger)

			foundLogger := getLogger(c)

			Convey("getLogger should return the logger", func() {
				So(foundLogger, ShouldNotBeNil)
				So(foundLogger, ShouldEqual, logger)
			})
		})

		Convey("When that context does not have a token", func() {
			c := e.NewContext(req, rec)
			foundToken := getToken(c)

			Convey("getToken should return nil", func() {
				So(foundToken, ShouldBeNil)
			})

		})

		Convey("When that context has a token", func() {
			c := e.NewContext(req, rec)
			token := &jwt.Token{}
			c.Set(tokenKey, token)

			foundToken := getToken(c)

			Convey("getToken should return that token", func() {
				So(foundToken, ShouldEqual, token)
			})

		})

		Convey("When that context does not have a dbSession", func() {
			c := e.NewContext(req, rec)
			foundDbSession := getAppDb(c)

			Convey("getAppDb should return nil", func() {
				So(foundDbSession, ShouldBeNil)
			})

		})

		Convey("When that context has a dbSession", func() {
			c := e.NewContext(req, rec)

			db := data.AppDb{SaltHash: "something cool and diff"}
			c.Set(dbKey, &db)

			foundDbSession := getAppDb(c)

			Convey("getAppDb should return the db", func() {
				So(foundDbSession, ShouldPointTo, &db)
			})

		})
	})

}
