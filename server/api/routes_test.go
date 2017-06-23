package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"world-backup/server/conf"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestSetUpRoutes(t *testing.T) {
	Convey("Given api.SetUpRoutes", t, func() {
		echoMock := new(EchoServerMock)

		api := &API{
			config: &conf.Config{Port: 7630},
			Server: echoMock,
			log:    logrus.WithField("test", "TestSetUpRoutes"),
		}

		groupMock := new(EchoGroupMock)

		echoMock.On("Use", mock.Anything, mock.Anything).Times(3)
		echoMock.On("GET", "*", mock.Anything, mock.Anything).Once()

		echoMock.On("Group", "/api", mock.Anything, mock.Anything).Return(groupMock)

		groupMock.On("GET", "/folders", mock.Anything, mock.Anything).Once()
		groupMock.On("GET", "/folders/:id/worlds", mock.Anything, mock.Anything).Once()
		groupMock.On("DELETE", "/folders/:id/worlds/:wid/backups/:bid", mock.Anything, mock.Anything).Once()

		routes := []echo.Route{
			{Path: "/something", Method: "put"},
		}
		echoMock.On("Routes").Return(routes)

		Convey("It should setup the Server Server properly", func() {
			api.SetUpRoutes()

			echoMock.AssertExpectations(t)
			groupMock.AssertExpectations(t)

			So(echoMock.Calls[3].Arguments.Get(0), ShouldEqual, "*")
			So(echoMock.Calls[3].Arguments.Get(1), ShouldEqual, api.index)

			So(groupMock.Calls[0].Arguments.Get(0), ShouldEqual, "/folders")
			So(groupMock.Calls[0].Arguments.Get(1), ShouldEqual, api.getFolders)

			So(groupMock.Calls[1].Arguments.Get(0), ShouldEqual, "/folders/:id/worlds")
			So(groupMock.Calls[1].Arguments.Get(1), ShouldEqual, api.getWorlds)

			So(groupMock.Calls[2].Arguments.Get(0), ShouldEqual, "/folders/:id/worlds/backups/:bid")
			So(groupMock.Calls[2].Arguments.Get(1), ShouldEqual, api.deleteWorldBackup)
		})

	})
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
