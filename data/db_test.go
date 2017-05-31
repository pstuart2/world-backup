package data

import (
	"testing"

	"time"

	"os"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
	"github.com/ventu-io/go-shortid"
)

var db *Db
var af afero.Afero

func TestMain(m *testing.M) {
	dbName := "test_" + shortid.MustGenerate() + ".json"

	af = afero.Afero{Fs: afero.NewOsFs()}
	db = Open(dbName, af)

	code := m.Run()

	af.Remove(dbName)
	os.Exit(code)
}

func TestOpen(t *testing.T) {

	Convey("Given a name and filesystem", t, func() {
		now := time.Now()
		before := now.AddDate(0, 0, -1)

		oldGetNow := getNow
		getNow = func() time.Time { return before }
		defer func() { getNow = oldGetNow }()

		dbName := "somd.json"

		fs := afero.Afero{Fs: afero.NewOsFs()}
		localDb := Open(dbName, fs)

		defer func() {
			fs.Remove(dbName)
		}()

		getNow = func() time.Time { return now }

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
				db2 := Open(dbName, fs)

				So(db2, ShouldNotBeNil)
				So(db2.data.LastSave.UnixNano(), ShouldEqual, now.UnixNano())

				Convey("close should also save the localDb", func() {
					next := now.AddDate(0, 0, 1)
					getNow = func() time.Time { return next }

					db2.Close()

					db3 := Open(dbName, fs)
					So(db3.data.LastSave.UnixNano(), ShouldEqual, next.UnixNano())
				})
			})

		})

	})

}
