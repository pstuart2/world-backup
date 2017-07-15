package fs

import (
	"testing"

	"errors"

	"path"

	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFileSystem_CreateBackup(t *testing.T) {
	Convey("Given an IBackupFs", t, func() {
		folderPath := "/path/to/folder"
		worldName := "backup-world-001"
		backupDir := "/path/to/backups"
		backupName := "ThisBeTheBackup.zip"

		fsMock := new(IBackupFsMock)

		log := logrus.WithField("test", "fs")

		Convey("When we fail to get current working directory", func() {
			fsMock.On("Getwd").Return("", errors.New("No!"))

			err := CreateBackup(fsMock, log, folderPath, worldName, backupDir, backupName)

			Convey("Then it should return an error", func() {
				fsMock.AssertExpectations(t)

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "No!")
			})
		})

		Convey("When the backup succeeds", func() {

			fsMock.On("Getwd").Return("/app/dir", nil)
			fsMock.On("Chdir", folderPath).Return(nil)
			fsMock.On("Chdir", "/app/dir").Return(nil)
			fsMock.On("Zip", worldName, path.Join(backupDir, backupName)).Return(nil)

			err := CreateBackup(fsMock, log, folderPath, worldName, backupDir, backupName)

			fsMock.AssertExpectations(t)

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the backup fails", func() {
			fsMock.On("Getwd").Return("/app/dir", nil)
			fsMock.On("Chdir", folderPath).Return(nil)
			fsMock.On("Chdir", "/app/dir").Return(nil)
			fsMock.On("Zip", worldName, path.Join(backupDir, backupName)).Return(errors.New("Didn't work!"))

			err := CreateBackup(fsMock, log, folderPath, worldName, backupDir, backupName)

			fsMock.AssertExpectations(t)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Didn't work!")
			})
		})
	})
}

func TestFileSystem_CleanName(t *testing.T) {
	Convey("Given a name", t, func() {
		Convey("It should replace any non-alphanumeric with _", func() {
			badname := "! This # Would be a / not so good ^name"
			result := CleanName(badname)
			So(result, ShouldEqual, "_This_Would_be_a_not_so_good_name")
		})
	})
}
