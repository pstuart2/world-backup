package api

import (
	"errors"
	"net/http"
	"testing"

	"awana-app/server/conf"

	"strings"

	"net/http/httptest"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2"
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
		echoMock.On("GET", "*", mock.Anything).Once()
		echoMock.On("POST", "/auth/signup", mock.Anything, mock.Anything).Once()
		echoMock.On("POST", "/auth/login", mock.Anything, mock.Anything).Once()

		echoMock.On("Group", "/api", mock.Anything).Return(groupMock)

		groupMock.On("GET", "/authenticate", mock.Anything, mock.Anything).Once()

		groupMock.On("GET", "/clubTypes", mock.Anything, mock.Anything).Once()
		groupMock.On("GET", "/awards/:clubTypeId", mock.Anything, mock.Anything).Once()
		groupMock.On("POST", "/clubs", mock.Anything, mock.Anything).Once()
		groupMock.On("GET", "/clubs", mock.Anything, mock.Anything).Once()
		groupMock.On("POST", "/books", mock.Anything, mock.Anything).Once()
		groupMock.On("PATCH", "/books/:id", mock.Anything, mock.Anything).Once()
		groupMock.On("GET", "/books", mock.Anything, mock.Anything).Once()
		groupMock.On("GET", "/books/:id", mock.Anything, mock.Anything).Once()
		groupMock.On("POST", "/books/:bookId/chapters", mock.Anything, mock.Anything).Once()
		groupMock.On("POST", "/books/:bookId/chapters/:chapterId", mock.Anything, mock.Anything).Once()
		groupMock.On("DELETE", "/books/:bookId/chapters/:chapterId", mock.Anything, mock.Anything).Once()

		groupMock.On("GET", "*", mock.Anything).Once()

		routes := []echo.Route{
			{Path: "/something", Method: "put"},
		}
		echoMock.On("Routes").Return(routes)

		Convey("It should setup the Server Server properly", func() {
			api.SetUpRoutes()

			echoMock.AssertExpectations(t)

			So(echoMock.Calls[2].Arguments.Get(0), ShouldEqual, "*")
			So(echoMock.Calls[2].Arguments.Get(1), ShouldEqual, api.index)
			So(echoMock.Calls[3].Arguments.Get(0), ShouldEqual, "/auth/signup")
			So(echoMock.Calls[3].Arguments.Get(1), ShouldEqual, api.signUp)
			So(echoMock.Calls[4].Arguments.Get(0), ShouldEqual, "/auth/login")
			So(echoMock.Calls[4].Arguments.Get(1), ShouldEqual, api.login)
			So(echoMock.Calls[5].Arguments.Get(0), ShouldEqual, "/api")
			So(echoMock.Calls[5].Arguments.Get(1), ShouldNotBeNil)

			groupMock.AssertExpectations(t)

			So(groupMock.Calls[0].Arguments.Get(0), ShouldEqual, "/authenticate")
			So(groupMock.Calls[0].Arguments.Get(1), ShouldEqual, api.authenticate)
			So(groupMock.Calls[1].Arguments.Get(0), ShouldEqual, "/clubTypes")
			So(groupMock.Calls[1].Arguments.Get(1), ShouldEqual, api.clubTypes)
			So(groupMock.Calls[2].Arguments.Get(0), ShouldEqual, "/awards/:clubTypeId")
			So(groupMock.Calls[2].Arguments.Get(1), ShouldEqual, api.awards)

			So(groupMock.Calls[3].Arguments.Get(0), ShouldEqual, "/clubs")
			So(groupMock.Calls[3].Arguments.Get(1), ShouldEqual, api.createClub)
			So(groupMock.Calls[4].Arguments.Get(0), ShouldEqual, "/clubs")
			So(groupMock.Calls[4].Arguments.Get(1), ShouldEqual, api.clubs)

			So(groupMock.Calls[5].Arguments.Get(0), ShouldEqual, "/books")
			So(groupMock.Calls[5].Arguments.Get(1), ShouldEqual, api.createBook)
			So(groupMock.Calls[6].Arguments.Get(0), ShouldEqual, "/books/:id")
			So(groupMock.Calls[6].Arguments.Get(1), ShouldEqual, api.updateBook)
			So(groupMock.Calls[7].Arguments.Get(0), ShouldEqual, "/books")
			So(groupMock.Calls[7].Arguments.Get(1), ShouldEqual, api.books)
			So(groupMock.Calls[8].Arguments.Get(0), ShouldEqual, "/books/:id")
			So(groupMock.Calls[8].Arguments.Get(1), ShouldEqual, api.book)

			So(groupMock.Calls[9].Arguments.Get(0), ShouldEqual, "/books/:bookId/chapters")
			So(groupMock.Calls[9].Arguments.Get(1), ShouldEqual, api.addChapter)
			So(groupMock.Calls[10].Arguments.Get(0), ShouldEqual, "/books/:bookId/chapters/:chapterId")
			So(groupMock.Calls[10].Arguments.Get(1), ShouldEqual, api.updateChapter)
			So(groupMock.Calls[11].Arguments.Get(0), ShouldEqual, "/books/:bookId/chapters/:chapterId")
			So(groupMock.Calls[11].Arguments.Get(1), ShouldEqual, api.deleteChapter)
		})

	})

}

func TestSetupRequest(t *testing.T) {
	Convey("Given an API object", t, func() {
		echoMock := new(EchoServerMock)

		session, _ := mgo.Dial("mongodb://localhost/awana-app-test")
		defer session.Close()

		api := &API{
			config:    &conf.Config{Port: 7630},
			Server:    echoMock,
			dbSession: session,
			log:       logrus.WithField("test", "TestSetupRequest"),
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
			So(c.Get(dbKey), ShouldNotBeNil)

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
			So(c.Get(dbKey), ShouldNotBeNil)

			handlerMock.AssertExpectations(t)

			Convey("The error should not be set", func() {
				So(handlerMock.Calls[0].Arguments[0], ShouldEqual, c)
				So(handlerError, ShouldEqual, expectedError)
				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
