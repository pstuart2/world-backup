package data

import (
	"testing"

	"time"

	"os"

	"errors"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/ventu-io/go-shortid"
)

var db *Db
var af afero.Afero

func TestMain(m *testing.M) {
	dbName := "test_" + shortid.MustGenerate() + ".json"

	af = afero.Afero{Fs: afero.NewOsFs()}
	db, _ = Open(dbName, af)

	code := m.Run()

	af.Remove(dbName)
	os.Exit(code)
}

func TestOpenAndSave(t *testing.T) {

	Convey("Given a name and filesystem", t, func() {
		now := time.Now()
		before := now.AddDate(0, 0, -1)

		oldGetNow := getNow
		getNow = func() time.Time { return before }
		defer func() { getNow = oldGetNow }()

		dbName := "somd.json"

		fs := afero.Afero{Fs: afero.NewOsFs()}
		localDb, ldbErr := Open(dbName, fs)

		defer func() {
			fs.Remove(dbName)
		}()

		getNow = func() time.Time { return now }

		So(ldbErr, ShouldBeNil)
		So(localDb, ShouldNotBeNil)
		So(localDb.name, ShouldEqual, dbName)
		So(localDb.data.CreatedAt.UnixNano(), ShouldEqual, before.UnixNano())
		So(localDb.data.LastSave.UnixNano(), ShouldEqual, -6795364578871345152)

		Convey("It should add save it to the localDb", func() {
			saveError := localDb.Save()
			So(saveError, ShouldBeNil)
			So(localDb.data.CreatedAt.UnixNano(), ShouldEqual, before.UnixNano())
			So(localDb.data.LastSave.UnixNano(), ShouldEqual, now.UnixNano())

			Convey("and be able to read it back", func() {
				db2, _ := Open(dbName, fs)

				So(db2, ShouldNotBeNil)
				So(db2.data.LastSave.UnixNano(), ShouldEqual, now.UnixNano())

				Convey("close should also save the localDb", func() {
					next := now.AddDate(0, 0, 1)
					getNow = func() time.Time { return next }

					db2.Close()

					db3, _ := Open(dbName, fs)
					So(db3.data.LastSave.UnixNano(), ShouldEqual, next.UnixNano())
				})
			})

		})

	})

	Convey("Given there are some file system errors", t, func() {
		fsMock := new(IDbFileSystemMock)

		Convey("When the exists check errors", func() {
			fsMock.On("Exists", "aName").Return(false, errors.New("What name?"))

			_, err := Open("aName", fsMock)

			Convey("It should return the error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "What name?")
			})
		})

		Convey("When we fail to read a file", func() {
			fsMock.On("Exists", "aName").Return(true, nil)
			fsMock.On("ReadFile", "aName").Return(nil, errors.New("failed to read"))

			_, err := Open("aName", fsMock)

			Convey("It should return the error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "failed to read")
			})
		})

		Convey("When the file is junk", func() {
			fsMock.On("Exists", "aName").Return(true, nil)
			fsMock.On("ReadFile", "aName").Return([]byte("Junk!"), nil)

			_, err := Open("aName", fsMock)

			Convey("It should return the error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid character 'J' looking for beginning of value")
			})
		})

		Convey("When we fail to write the file", func() {
			fsMock.On("WriteFile", "Wow fake", mock.Anything, mock.Anything).Return(errors.New("NOOO!"))
			db := Db{
				fs:   fsMock,
				name: "Wow fake",
			}
			err := db.Save()

			Convey("It should return the error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "NOOO!")
			})
		})
	})

}
