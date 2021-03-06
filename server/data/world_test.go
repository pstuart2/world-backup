package data

import (
	"testing"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWorld_AddBackup(t *testing.T) {

	Convey("Given a world", t, func() {
		oldGetId := getId
		idCounter := 0
		getId = func() string { idCounter++; return fmt.Sprintf("ID:%d", idCounter) }
		defer func() { getId = oldGetId }()

		world := World{Id: "C00L"}

		Convey("When we add a backup", func() {
			world.AddBackup("Backup001.zip")
			world.AddBackup("Backup002.zip")
			world.AddBackup("Backup003.zip")

			Convey("It should add the backup the end of backups", func() {
				So(len(world.Backups), ShouldEqual, 3)
				So(world.Backups[0].Id, ShouldEqual, "ID:1")
				So(world.Backups[0].Name, ShouldEqual, "Backup001.zip")
				So(world.Backups[1].Id, ShouldEqual, "ID:2")
				So(world.Backups[1].Name, ShouldEqual, "Backup002.zip")
				So(world.Backups[2].Id, ShouldEqual, "ID:3")
				So(world.Backups[2].Name, ShouldEqual, "Backup003.zip")

				So(world.Backups[2].CreatedAt.After(world.Backups[0].CreatedAt), ShouldBeTrue)
			})

		})

	})

}

func TestWorld_LastBackup(t *testing.T) {
	Convey("Given a world", t, func() {
		world := World{}

		Convey("With no backups", func() {
			Convey("Then it should return 0 time", func() {
				t := world.LastBackupTime()
				So(t.IsZero(), ShouldBeTrue)
			})
		})

		Convey("With backups", func() {
			b1 := world.AddBackup("b01.zip")
			world.AddBackup("b02.zip")
			b3 := world.AddBackup("b03.zip")

			Convey("It should return the last back up date", func() {
				t := world.LastBackupTime()
				So(t.Equal(b1.CreatedAt), ShouldBeFalse)
				So(t.Equal(b3.CreatedAt), ShouldBeTrue)
			})
		})
	})
}

func TestWorld_RemoveBackup(t *testing.T) {
	Convey("Given a world with backups", t, func() {
		world := World{}
		world.AddBackup("Backup 001")
		world.AddBackup("Backup 002")
		b3 := world.AddBackup("Backup 003")
		world.AddBackup("Backup 004")
		world.AddBackup("Backup 005")

		idToRemove := b3.Id

		Convey("When a backup is removed", func() {
			world.RemoveBackup(idToRemove)

			Convey("It should remove the backup and keep the others", func() {
				So(len(world.Backups), ShouldEqual, 4)
				So(world.findBackupIndex(idToRemove), ShouldEqual, -1)
			})
		})
	})
}

func TestWorld_GetBackup(t *testing.T) {
	Convey("Given a world with backups", t, func() {
		world := World{}
		world.AddBackup("Backup 001")
		world.AddBackup("Backup 002")
		b3 := world.AddBackup("Backup 003")
		b4 := world.AddBackup("Backup 004")
		world.AddBackup("Backup 005")

		Convey("When GetBackup is called with a valid id", func() {
			r1 := world.GetBackup(b3.Id)
			r2 := world.GetBackup(b4.Id)

			Convey("It should return the backup", func() {
				So(r1.Id, ShouldEqual, b3.Id)
				So(r2.Id, ShouldEqual, b4.Id)
			})
		})

		Convey("When GetBackup is called with an invalid id", func() {
			r3 := world.GetBackup("nonon")

			Convey("It should return nil", func() {
				So(r3, ShouldBeNil)
			})
		})
	})
}
