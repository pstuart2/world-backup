package api

import (
	"testing"

	"world-backup/conf"

	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
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

		Convey("It should return a new api object", func() {
			api := NewAPI(log, &conf)

			So(api, ShouldNotBeNil)
			So(api.config, ShouldEqual, &conf)
			So(api.Server, ShouldNotBeNil)
		})

	})

}
