package data

import (
	"testing"

	"time"

	"fmt"

	"path"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFolder_AddWorld(t *testing.T) {
	Convey("Given a folder", t, func() {
		oldGetId := getId
		idCounter := 0
		getId = func() string { idCounter++; return fmt.Sprintf("WID:%d", idCounter) }
		defer func() { getId = oldGetId }()

		folder := Folder{Id: "C00L", Path: "/some/path/for"}

		Convey("When we add a backup", func() {
			w1 := folder.AddWorld("World 01 Spot")
			w2 := folder.AddWorld("World 02 Spot")
			w3 := folder.AddWorld("World 03 Spot")

			So(w1.Id, ShouldEqual, "WID:1")
			So(w1.FullPath, ShouldEqual, path.Join(folder.Path, w1.Name))
			So(w2.Id, ShouldEqual, "WID:2")
			So(w2.FullPath, ShouldEqual, path.Join(folder.Path, w2.Name))
			So(w3.Id, ShouldEqual, "WID:3")
			So(w3.FullPath, ShouldEqual, path.Join(folder.Path, w3.Name))

			Convey("It should add the backup the end of backups", func() {
				So(len(folder.Worlds), ShouldEqual, 3)
				So(folder.Worlds[0].Id, ShouldEqual, w1.Id)
				So(folder.Worlds[0].Name, ShouldEqual, "World 01 Spot")
				So(folder.Worlds[1].Id, ShouldEqual, w2.Id)
				So(folder.Worlds[1].Name, ShouldEqual, "World 02 Spot")
				So(folder.Worlds[2].Id, ShouldEqual, w3.Id)
				So(folder.Worlds[2].Name, ShouldEqual, "World 03 Spot")

				So(folder.Worlds[2].CreatedAt.After(folder.Worlds[0].CreatedAt), ShouldBeTrue)
			})

		})

	})
}

func TestFolder_GetWorldByName(t *testing.T) {
	Convey("Given a folder with multiple worlds", t, func() {
		folder := Folder{}
		w1 := folder.AddWorld("World Numeral 01")
		folder.AddWorld("World Numeral 02")
		w3 := folder.AddWorld("World Numeral 03")
		folder.AddWorld("World Numeral 04")

		Convey("It should be able to find the worlds by name", func() {
			r1 := folder.GetWorldByName("World Numeral 01")
			r2 := folder.GetWorldByName("World Numeral 03")
			r3 := folder.GetWorldByName("Not here")

			So(r1.Id, ShouldEqual, w1.Id)
			So(r2.Id, ShouldEqual, w3.Id)
			So(r3, ShouldBeNil)
		})
	})
}

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
					db2 := Open(db.name, af)

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
			r1 := db.GetFolderByPath("/some/cool2/place")
			r2 := db.GetFolderByPath("/some/cool4/place")
			r3 := db.GetFolderByPath("/some/not/here")

			So(r1.Id, ShouldEqual, f2.Id)
			So(r2.Id, ShouldEqual, f4.Id)
			So(r3, ShouldBeNil)
		})
	})
}
