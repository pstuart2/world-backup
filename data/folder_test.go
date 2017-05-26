package data

import (
	"testing"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFolder_AddFolder(t *testing.T) {
	now := time.Now()

	oldGetNow := getNow
	getNow = func() time.Time { return now }
	defer func() { getNow = oldGetNow }()

	oldGetId := getId
	getId = func() string { return "TheCoolID" }
	defer func() { getId = oldGetId }()

	Convey("Given a folder", t, func() {

		Convey("It should add it to the db and return it", func() {

			f1 := db.AddFolder("/some/cool/place")
			So(len(db.data.Folders), ShouldEqual, 1)

			So(f1, ShouldNotBeNil)
			So(f1.Id, ShouldEqual, "TheCoolID")
			So(f1.Path, ShouldEqual, "/some/cool/place")
			So(f1.CreatedAt.UnixNano(), ShouldEqual, now.UnixNano())
			So(f1.ModifiedAt.UnixNano(), ShouldEqual, now.UnixNano())

			Convey("It should be able to add more", func() {
				getId = func() string { return "TheCoolID2" }

				f2 := db.AddFolder("another cool place")
				So(f2.Path, ShouldEqual, "another cool place")
				So(len(db.data.Folders), ShouldEqual, 2)

				Convey("and be able to read it back", func() {
					db.Save()
					db2 := Open(db.name, fs)

					So(len(db2.data.Folders), ShouldEqual, 2)
					So(db2.data.Folders[0].Id, ShouldEqual, "TheCoolID")
					So(db2.data.Folders[0].Path, ShouldEqual, "/some/cool/place")
					So(db2.data.Folders[1].Id, ShouldEqual, "TheCoolID2")
					So(db2.data.Folders[1].Path, ShouldEqual, "another cool place")

				})
			})

		})

	})

}

func TestFolder_FolderByPath(t *testing.T) {
	Convey("Given a list of folders", t, func() {
		db.AddFolder("/some/cool1/place")
		f2 := db.AddFolder("/some/cool2/place")
		db.AddFolder("/some/cool3/place")
		f4 := db.AddFolder("/some/cool4/place")
		db.AddFolder("/some/cool5/place")

		Convey("It should return the correct folder when asked", func() {
			r1 := db.FolderByPath("/some/cool2/place")
			r2 := db.FolderByPath("/some/cool4/place")

			So(r1.Id, ShouldEqual, f2.Id)
			So(r2.Id, ShouldEqual, f4.Id)
		})
	})
}
