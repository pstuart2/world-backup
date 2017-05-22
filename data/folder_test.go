package data

import (
	"testing"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddFolder(t *testing.T) {
	now := time.Now()

	oldGetNow := getNow
	getNow = func() time.Time { return now }
	defer func() { getNow = oldGetNow }()

	Convey("Given a folder", t, func() {

		Convey("It should add it to the db and return it", func() {

			f1 := db.AddFolder("/some/cool/place")
			So(len(db.data.Folders), ShouldEqual, 1)

			So(f1, ShouldNotBeNil)
			So(f1.Path, ShouldEqual, "/some/cool/place")
			So(f1.CreatedAt.UnixNano(), ShouldEqual, now.UnixNano())
			So(f1.ModifiedAt.UnixNano(), ShouldEqual, now.UnixNano())

			Convey("It should be able to add more", func() {
				f2 := db.AddFolder("another cool place")
				So(f2.Path, ShouldEqual, "another cool place")
				So(len(db.data.Folders), ShouldEqual, 2)

				Convey("and be able to read it back", func() {
					db.Save()
					db2 := Open(db.name, fs)

					So(len(db2.data.Folders), ShouldEqual, 2)
					So(db2.data.Folders[0].Path, ShouldEqual, "/some/cool/place")
					So(db2.data.Folders[1].Path, ShouldEqual, "another cool place")

				})
			})

		})

	})

}
