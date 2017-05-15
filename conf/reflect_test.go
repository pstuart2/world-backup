package conf

import (
	"testing"

	"reflect"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestSimpleValues(t *testing.T) {

	Convey("Given some simple values", t, func() {
		c := struct {
			Simple string `json:"simple"`
		}{}

		viper.SetDefault("simple", "i am a simple string")

		Convey("It should be able set and get the values out", func() {
			err := recursivelySet(reflect.ValueOf(&c), "")

			So(err, ShouldBeNil)
			So(c.Simple, ShouldEqual, "i am a simple string")
		})

	})

}

func TestNewValues(t *testing.T) {

	Convey("Given a nested structure", t, func() {

		c := struct {
			Simple string `json:"simple"`
			Nested struct {
				BoolVal   bool   `json:"bool"`
				StringVal string `json:"string"`
				NumberVal int    `json:"number"`
			} `json:"nested"`
		}{}

		viper.SetDefault("simple", "simple")
		viper.SetDefault("nested.bool", true)
		viper.SetDefault("nested.string", "i am a simple string")
		viper.SetDefault("nested.number", 4)

		Convey("It should be able to set and get the values out", func() {
			err := recursivelySet(reflect.ValueOf(&c), "")

			So(err, ShouldBeNil)
			So(c.Simple, ShouldEqual, "simple")
			So(c.Nested.NumberVal, ShouldEqual, 4)
			So(c.Nested.StringVal, ShouldEqual, "i am a simple string")
			So(c.Nested.BoolVal, ShouldEqual, true)
		})

	})

}
