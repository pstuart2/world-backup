package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"world-backup/conf"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestSetUpRoutes(t *testing.T) {

}

func TestSetupRequest(t *testing.T) {
	Convey("Given an API object", t, func() {
		echoMock := new(EchoServerMock)

		api := &API{
			config: &conf.Config{Port: 7630},
			Server: echoMock,
			log:    logrus.WithField("test", "TestSetupRequest"),
		}

		Convey("When the request is successful", func() {
			handlerMock := new(HandlerMock)
			handlerMock.On("Handler", mock.Anything).Return(nil)

			e := echo.New()
			req, _ := http.NewRequest(echo.GET, "/users", strings.NewReader(""))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handlerError := api.setupRequest(handlerMock.Handler)(c)

			So(c.Get(loggerKey), ShouldNotBeNil)

			handlerMock.AssertExpectations(t)

			Convey("The error should not be set", func() {
				So(handlerMock.Calls[0].Arguments[0], ShouldEqual, c)
				So(handlerError, ShouldBeNil)
				So(rec.Code, ShouldEqual, http.StatusOK)
			})
		})

		Convey("When the request has an error", func() {
			expectedError := errors.New("Something went wrong")

			handlerMock := new(HandlerMock)
			handlerMock.On("Handler", mock.Anything).Return(expectedError)

			e := echo.New()
			req, _ := http.NewRequest(echo.GET, "/users", strings.NewReader(""))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handlerError := api.setupRequest(handlerMock.Handler)(c)

			So(c.Get(loggerKey), ShouldNotBeNil)

			handlerMock.AssertExpectations(t)

			Convey("The error should not be set", func() {
				So(handlerMock.Calls[0].Arguments[0], ShouldEqual, c)
				So(handlerError, ShouldEqual, expectedError)
				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
