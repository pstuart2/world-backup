package api

import (
	"awana-app/server/conf"
	"testing"

	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
)

func TestStart(t *testing.T) {

	Convey("Given an api object", t, func() {
		echoMock := new(EchoServerMock)

		api := &API{
			config: &conf.Config{Port: 7630},
			Server: echoMock,
		}

		Convey("It should start the Server on the port from the config", func() {
			echoMock.On("Start", ":7630").Return(nil)

			api.Start()

			echoMock.AssertExpectations(t)
		})

	})

}

func TestNewAPI(t *testing.T) {

	Convey("Given a logger, config and an db", t, func() {
		conf := conf.Config{}
		log := logrus.WithField("test", "TestNewApi")
		session := mgo.Session{}

		Convey("It should return a new api object", func() {
			api := NewAPI(log, &conf, &session)

			So(api, ShouldNotBeNil)
			So(api.config, ShouldEqual, &conf)
			So(api.dbSession, ShouldEqual, &session)
			So(api.Server, ShouldNotBeNil)
		})

	})

}
