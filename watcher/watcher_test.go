package watcher

import (
	"testing"
	"world-backup/conf"

	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"time"
	"errors"
)

func TestWatcher_NewWatcher(t *testing.T) {

	Convey("Given the correct input", t, func() {
		config := conf.Config{}
		log := logrus.WithField("test", "watcher")
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		dbMock := new(IDbMock)
		zipMock := new(IArchiverMock)

		Convey("It should return a new watcher", func() {
			w := NewWatcher(log, &config, fs, dbMock, zipMock)

			So(w, ShouldNotBeNil)
			So(w.config, ShouldEqual, &config)
			So(w.db, ShouldEqual, dbMock)
			So(w.zip, ShouldEqual, zipMock)
		})

	})

}

func TestWatcher_Check(t *testing.T) {
	Convey("Given a watcher and directories that do not exist", t, func() {
		config := conf.Config{
			WatchDir:  "/home/world",
			BackupDir: "/home/backup",
		}
		log := logrus.WithField("test", "watcher")
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		dbMock := new(IDbMock)
		zipMock := new(IArchiverMock)

		w := NewWatcher(log, &config, fs, dbMock, zipMock)

		Convey("When there are no sub directories", func() {
			w.check()

			Convey("It should not create any backups", func() {
				So(len(zipMock.Calls), ShouldEqual, 0)
			})
		})
	})

	Convey("Given a watcher and a valid watch directory", t, func() {
		now := time.Unix(1495807405, 0)
		oldGetNow := getNow
		getNow = func() time.Time { return now }
		defer func() { getNow = oldGetNow }()

		config := conf.Config{
			WatchDir:  "/home/world",
			BackupDir: "/home/backup",
		}
		log := logrus.WithField("test", "watcher")
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		dbMock := new(IDbMock)
		zipMock := new(IArchiverMock)

		fs.MkdirAll("/home/world", 0755)
		fs.MkdirAll("/home/backup", 0755)

		w := NewWatcher(log, &config, fs, dbMock, zipMock)

		Convey("When there are no sub directories", func() {
			w.check()

			Convey("It should not create any backups", func() {
				So(len(zipMock.Calls), ShouldEqual, 0)
			})
		})

		Convey("When there are directories", func() {
			fs.MkdirAll("/home/world/World one", 0755)
			f1Err := fs.WriteFile("/home/world/World one/f1.mc", []byte("file 1"), 0644)
			So(f1Err, ShouldBeNil)

			fs.MkdirAll("/home/world/World two", 0755)
			f2Err := fs.WriteFile("/home/world/World two/f2.mc", []byte("file 2"), 0644)
			So(f2Err, ShouldBeNil)

			zipMock.On("Make", mock.Anything, mock.Anything).Times(2).Return(nil)

			w.check()

			Convey("It should create the corresponding zip backups", func() {
				zipMock.AssertExpectations(t)

				c1 := zipMock.Calls[0];
				So(c1.Arguments.Get(0), ShouldEqual, "/home/backup/World one-20170526T090325.zip")
				So(len(c1.Arguments.Get(1).([]string)), ShouldEqual, 1)
				So(c1.Arguments.Get(1).([]string)[0], ShouldEqual, "/home/world/World one")

				c2 := zipMock.Calls[1];
				So(c2.Arguments.Get(0), ShouldEqual, "/home/backup/World two-20170526T090325.zip")
				So(len(c2.Arguments.Get(1).([]string)), ShouldEqual, 1)
				So(c2.Arguments.Get(1).([]string)[0], ShouldEqual, "/home/world/World two")
			})
		})

		Convey("When there is an error creating the zip", func() {
			fs.MkdirAll("/home/world/World one", 0755)
			f1Err := fs.WriteFile("/home/world/World one/f1.mc", []byte("file 1"), 0644)
			So(f1Err, ShouldBeNil)

			zipMock.On("Make", mock.Anything, mock.Anything).Return(errors.New("Oops!"))

			Convey("It should continue", func() {
				w.check()

				zipMock.AssertExpectations(t)
			})
		})

	})

}
