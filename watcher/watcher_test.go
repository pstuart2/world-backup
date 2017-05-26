package watcher

import (
	"testing"
	"world-backup/conf"

	"errors"
	"time"
	"world-backup/data"

	"os"

	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
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

func TestWatcher_Start(t *testing.T) {
	Convey("Given no directories to watch", t, func() {
		config := conf.Config{
			WatchDirs: []string{"/home/world", "/another/one"},
		}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)
		zipMock := new(IArchiverMock)

		w := NewWatcher(log, &config, fsMock, dbMock, zipMock)

		wasChecked := false
		oldCheck := check
		check = func(w *Watcher) { wasChecked = true }
		defer func() { check = oldCheck }()

		wasWatched := false
		oldWatch := watch
		watch = func(w *Watcher, stop chan bool) {
			wasWatched = true
		}
		defer func() { watch = oldWatch }()

		Convey("It should create the folders in the db that do not exist", func() {
			f1 := data.Folder{Id: "SomeId01", Path: "/home/world"}
			f2 := data.Folder{Id: "SomeId02", Path: "/another/one"}

			dbMock.On("FolderByPath", "/home/world").Return(&f1)
			dbMock.On("FolderByPath", "/another/one").Return(nil)
			dbMock.On("AddFolder", "/another/one").Return(&f2)
			dbMock.On("Save").Return(nil)

			w.Start()

			dbMock.AssertExpectations(t)

			Convey("and call check() and watch() at start", func() {
				So(wasChecked, ShouldBeTrue)

				<-time.After(time.Millisecond * 200)
				So(wasWatched, ShouldBeTrue)
			})
		})
	})

	Convey("Given directories to watch", t, func() {
		config := conf.Config{}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)
		zipMock := new(IArchiverMock)

		w := NewWatcher(log, &config, fsMock, dbMock, zipMock)

		Convey("It should not do anything", func() {
			w.Start()

			So(len(dbMock.Calls), ShouldEqual, 0)
		})
	})
}

func TestWatcher_Watch(t *testing.T) {
	Convey("Given a watcher", t, func() {
		config := conf.Config{
			CheckIntervalSeconds: 1,
		}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)
		zipMock := new(IArchiverMock)

		w := NewWatcher(log, &config, fsMock, dbMock, zipMock)

		checkCount := 0
		oldCheck := check
		check = func(w *Watcher) { checkCount++ }
		defer func() { check = oldCheck }()

		Convey("It should watch until stopped", func() {
			stopChannel := make(chan bool)

			go watch(w, stopChannel)
			<-time.After(time.Millisecond * 2200)
			stopChannel <- true

			So(checkCount, ShouldEqual, 2)
		})
	})
}

func TestWatcher_Check(t *testing.T) {
	Convey("Given a watcher to watch multiple directories", t, func() {
		config := conf.Config{
			WatchDirs: []string{"/home/paul", "/home/sydney", "/home/logan"},
		}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)
		zipMock := new(IArchiverMock)

		w := NewWatcher(log, &config, fsMock, dbMock, zipMock)

		folders := []*data.Folder{}
		oldCheckOneDir := checkOneDir
		checkOneDir = func(w *Watcher, f *data.Folder) {
			folders = append(folders, f)
		}
		defer func() { checkOneDir = oldCheckOneDir }()

		Convey("It should get the folder from the db", func() {
			f1 := data.Folder{Id: "01", Path: w.config.WatchDirs[0]}
			f2 := data.Folder{Id: "02", Path: w.config.WatchDirs[1]}
			f3 := data.Folder{Id: "03", Path: w.config.WatchDirs[2]}

			dbMock.On("FolderByPath", w.config.WatchDirs[0]).Return(&f1)
			dbMock.On("FolderByPath", w.config.WatchDirs[1]).Return(&f2)
			dbMock.On("FolderByPath", w.config.WatchDirs[2]).Return(&f3)

			Convey("and call checkOneDir for each", func() {
				check(w)

				dbMock.AssertExpectations(t)

				So(len(folders), ShouldEqual, 3)
				So(folders[0].Id, ShouldEqual, f1.Id)
				So(folders[1].Id, ShouldEqual, f2.Id)
				So(folders[2].Id, ShouldEqual, f3.Id)
			})
		})
	})
}

func TestWatcher_CheckOneDir(t *testing.T) {
	Convey("Given a watcher and directories", t, func() {
		config := conf.Config{
			BackupDir: "/home/backup",
		}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)
		zipMock := new(IArchiverMock)

		w := NewWatcher(log, &config, fsMock, dbMock, zipMock)

		Convey("When there is an error reading the directory", func() {
			fsMock.On("ReadDir", "/home/world").Return(nil, errors.New("No worky"))

			f := data.Folder{Path: "/home/world"}
			checkOneDir(w, &f)

			Convey("It should not create any backups", func() {
				fsMock.AssertExpectations(t)

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
			BackupDir: "/home/backup",
		}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)
		zipMock := new(IArchiverMock)

		w := NewWatcher(log, &config, fsMock, dbMock, zipMock)

		Convey("When there are no sub directories", func() {
			fsMock.On("ReadDir", "/home/world").Return([]os.FileInfo{}, nil)

			f := data.Folder{Path: "/home/world"}
			checkOneDir(w, &f)

			fsMock.AssertExpectations(t)

			Convey("It should not create any backups", func() {
				So(len(zipMock.Calls), ShouldEqual, 0)
			})
		})

		Convey("When there are directories", func() {
			dir1 := new(FileInfoMock)
			dir1.On("Name").Return("World one")
			dir1.On("IsDir").Return(true)
			dir1.On("ModTime").Return(now.Add(time.Second * -100))

			dir2 := new(FileInfoMock)
			dir2.On("Name").Return("World two")
			dir2.On("IsDir").Return(true)
			dir2.On("ModTime").Return(now.Add(time.Second * -100))

			fsMock.On("ReadDir", "/home/world").Return([]os.FileInfo{dir1, dir2}, nil)

			zipMock.On("Make", mock.Anything, mock.Anything).Times(2).Return(nil)

			f := data.Folder{Path: "/home/world"}
			checkOneDir(w, &f)

			Convey("It should create the corresponding zip backups", func() {
				fsMock.AssertExpectations(t)
				zipMock.AssertExpectations(t)

				c1 := zipMock.Calls[0]
				So(c1.Arguments.Get(0), ShouldEqual, "/home/backup/World one-20170526T090325.zip")
				So(len(c1.Arguments.Get(1).([]string)), ShouldEqual, 1)
				So(c1.Arguments.Get(1).([]string)[0], ShouldEqual, "/home/world/World one")

				c2 := zipMock.Calls[1]
				So(c2.Arguments.Get(0), ShouldEqual, "/home/backup/World two-20170526T090325.zip")
				So(len(c2.Arguments.Get(1).([]string)), ShouldEqual, 1)
				So(c2.Arguments.Get(1).([]string)[0], ShouldEqual, "/home/world/World two")
			})
		})

		Convey("When there is an error creating the zip", func() {
			dir1 := new(FileInfoMock)
			dir1.On("Name").Return("World one")
			dir1.On("IsDir").Return(true)
			dir1.On("ModTime").Return(now.Add(time.Second * -100))

			fsMock.On("ReadDir", "/home/world").Return([]os.FileInfo{dir1}, nil)

			zipMock.On("Make", mock.Anything, mock.Anything).Return(errors.New("Oops!"))

			Convey("It should continue", func() {
				f := data.Folder{Path: "/home/world"}
				checkOneDir(w, &f)

				fsMock.AssertExpectations(t)
				zipMock.AssertExpectations(t)
			})
		})

	})

}
