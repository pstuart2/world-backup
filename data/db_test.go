package data

import (
	"testing"

	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
)

var db *Db
var fs afero.Fs

//func TestMain(m *testing.M) {
//	dbName := "test_" + shortid.MustGenerate() + ".db"
//	db = Open(dbName)
//	defer db.Close()
//
//	code := m.Run()
//
//	fs = afero.NewOsFs()
//	fs.Remove(dbName)
//	os.Exit(code)
//}

func TestOpen(t *testing.T) {

	Convey("Given a name and filesystem", t, func() {
		now := time.Now()
		before := now.AddDate(0, 0, -1)

		oldGetNow := getNow
		getNow = func() time.Time { return before }
		defer func() { getNow = oldGetNow }()

		dbName := "somd.db"

		fs := afero.NewOsFs()
		db := Open(dbName, fs)

		defer func() {
			fs.Remove(dbName)
		}()

		getNow = func() time.Time { return now }

		So(db, ShouldNotBeNil)
		So(db.name, ShouldEqual, dbName)
		So(db.data.CreatedAt.UnixNano(), ShouldEqual, before.UnixNano())
		So(db.data.LastSave.UnixNano(), ShouldEqual, -6795364578871345152)

		Convey("It should add save it to the db", func() {
			saveError := db.Save()
			So(saveError, ShouldBeNil)
			So(db.data.CreatedAt.UnixNano(), ShouldEqual, before.UnixNano())
			So(db.data.LastSave.UnixNano(), ShouldEqual, now.UnixNano())

			Convey("and be able to read it back", func() {
				db2 := Open(dbName, fs)

				So(db2, ShouldNotBeNil)
				So(db2.data.LastSave.UnixNano(), ShouldEqual, now.UnixNano())
			})

		})

	})

}

func TestAddFolder(t *testing.T) {

	Convey("Given a folder", t, func() {

		Convey("It should add save it to the db", func() {

			Convey("and be able to read it back", func() {

			})

		})

	})

}
