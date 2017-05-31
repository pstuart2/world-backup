package watcher

import (
	"testing"
	"world-backup/conf"

	"errors"
	"time"
	"world-backup/data"

	"os"

	"fmt"

	"path"

	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
)

func TestWatcher_NewWatcher(t *testing.T) {

	Convey("Given the correct input", t, func() {
		config := conf.Config{}
		log := logrus.WithField("test", "watcher")
		fs := new(IFileSystemMock)
		dbMock := new(IDbMock)

		Convey("It should return a new watcher", func() {
			w := NewWatcher(log, &config, fs, dbMock)

			So(w, ShouldNotBeNil)
			So(w.config, ShouldEqual, &config)
			So(w.db, ShouldEqual, dbMock)
			So(w.fs, ShouldEqual, fs)
		})

	})

}

func TestWatcher_Start(t *testing.T) {
	Convey("Given directories to watch", t, func() {
		config := conf.Config{
			WatchDirs:     []string{"/home/world", "/another/one"},
			CheckInterval: "1s",
		}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)

		w := NewWatcher(log, &config, fsMock, dbMock)

		wasChecked := false
		oldCheck := check
		check = func(w *Watcher) { wasChecked = true }
		defer func() { check = oldCheck }()

		wasWatched := false
		oldWatch := watch
		watch = func(w *Watcher, stop chan bool, d time.Duration) {
			wasWatched = true
		}
		defer func() { watch = oldWatch }()

		Convey("It should create the folders in the db that do not exist", func() {
			f1 := data.Folder{Id: "SomeId01", Path: "/home/world"}
			f2 := data.Folder{Id: "SomeId02", Path: "/another/one"}

			dbMock.On("GetFolderByPath", "/home/world").Return(&f1)
			dbMock.On("GetFolderByPath", "/another/one").Return(nil)
			dbMock.On("AddFolder", "/another/one").Return(&f2)
			dbMock.On("Save").Return(nil)

			err := w.Start()

			So(err, ShouldBeNil)
			dbMock.AssertExpectations(t)

			Convey("and call check() and watch() at start", func() {
				So(wasChecked, ShouldBeTrue)

				<-time.After(time.Millisecond * 200)
				So(wasWatched, ShouldBeTrue)
			})
		})
	})

	Convey("Given no directories to watch", t, func() {
		config := conf.Config{}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)

		w := NewWatcher(log, &config, fsMock, dbMock)

		Convey("It should return NoWatchPathError", func() {
			err := w.Start()

			So(err, ShouldEqual, NoWatchPathError)
			So(len(dbMock.Calls), ShouldEqual, 0)
		})
	})

	Convey("Given an invalid watch interval", t, func() {
		config := conf.Config{
			WatchDirs:     []string{"/home/world", "/another/one"},
			CheckInterval: "33",
		}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)

		w := NewWatcher(log, &config, fsMock, dbMock)

		Convey("It should return InvalidCheckInterval", func() {
			err := w.Start()

			So(err, ShouldEqual, InvalidCheckInterval)
			So(len(dbMock.Calls), ShouldEqual, 0)
		})
	})
}

func TestWatcher_Watch(t *testing.T) {
	Convey("Given a watcher", t, func() {
		config := conf.Config{}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)

		w := NewWatcher(log, &config, fsMock, dbMock)

		checkCount := 0
		oldCheck := check
		check = func(w *Watcher) { checkCount++ }
		defer func() { check = oldCheck }()

		Convey("It should watch until stopped", func() {
			stopChannel := make(chan bool)

			go watch(w, stopChannel, time.Second*1)
			<-time.After(time.Millisecond * 2200)
			stopChannel <- true

			So(checkCount, ShouldEqual, 2)
		})
	})
}

func TestWatcher_Check(t *testing.T) {
	Convey("Given a watcher to watch multiple directories", t, func() {
		now := time.Now()
		oldGetNow := getNow
		getNow = func() time.Time { return now }
		defer func() { getNow = oldGetNow }()

		config := conf.Config{
			WatchDirs: []string{"/home/paul", "/home/sydney", "/home/logan"},
		}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)

		w := NewWatcher(log, &config, fsMock, dbMock)

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

			dbMock.On("GetFolderByPath", w.config.WatchDirs[0]).Return(&f1)
			dbMock.On("GetFolderByPath", w.config.WatchDirs[1]).Return(&f2)
			dbMock.On("GetFolderByPath", w.config.WatchDirs[2]).Return(&f3)
			dbMock.On("Save").Times(3).Return(nil)

			Convey("and call checkOneDir for each", func() {
				check(w)

				dbMock.AssertExpectations(t)

				So(len(folders), ShouldEqual, 3)
				So(folders[0].Id, ShouldEqual, f1.Id)
				So(folders[1].Id, ShouldEqual, f2.Id)
				So(folders[2].Id, ShouldEqual, f3.Id)

				Convey("and update the LastRun on each folder", func() {
					So(f1.LastRun.UnixNano(), ShouldEqual, now.UnixNano())
					So(f2.LastRun.UnixNano(), ShouldEqual, now.UnixNano())
					So(f3.LastRun.UnixNano(), ShouldEqual, now.UnixNano())
				})
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

		w := NewWatcher(log, &config, fsMock, dbMock)

		hasChangedFilesCallCount := 0
		oldHasChangedFiles := hasChangedFiles
		defer func() { hasChangedFiles = oldHasChangedFiles }()
		hasChangedFiles = func(log *logrus.Entry, fs IFileSystem, world *data.World) bool {
			hasChangedFilesCallCount++
			return false
		}

		Convey("When there is an error reading the directory", func() {
			fsMock.On("ReadDir", "/home/world").Return(nil, errors.New("No worky"))

			f := data.Folder{Path: "/home/world"}
			checkOneDir(w, &f)

			Convey("It should not create any backups", func() {
				fsMock.AssertExpectations(t)
				So(hasChangedFilesCallCount, ShouldEqual, 0)
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

		w := NewWatcher(log, &config, fsMock, dbMock)

		hasChangedFilesCallCount := 0
		oldHasChangedFiles := hasChangedFiles
		defer func() { hasChangedFiles = oldHasChangedFiles }()
		hasChangedFiles = func(log *logrus.Entry, fs IFileSystem, world *data.World) bool {
			hasChangedFilesCallCount++
			return (hasChangedFilesCallCount % 2) != 0
		}

		backupCreatedCallCount := 0
		var backedUpWorld *data.World
		oldCreateBackup := createBackup
		defer func() { createBackup = oldCreateBackup }()
		createBackup = func(w *Watcher, log *logrus.Entry, f *data.Folder, world *data.World) {
			backupCreatedCallCount++
			backedUpWorld = world
		}

		oldCheckPurgeBackup := checkPurgeBackup
		defer func() { checkPurgeBackup = oldCheckPurgeBackup }()
		checkPurgeBackupCallCount := 0
		checkPurgeBackup = func(w *Watcher, log *logrus.Entry, world *data.World) { checkPurgeBackupCallCount++ }

		Convey("When there are no worlds", func() {
			fsMock.On("ReadDir", "/home/world").Return([]os.FileInfo{}, nil)

			f := data.Folder{Path: "/home/world"}

			checkOneDir(w, &f)

			fsMock.AssertExpectations(t)

			Convey("It should not create any backups", func() {
				So(hasChangedFilesCallCount, ShouldEqual, 0)
				So(backupCreatedCallCount, ShouldEqual, 0)
			})
		})

		Convey("When there are worlds", func() {
			wDir1 := new(FileInfoMock)
			wDir1.On("Name").Return("World one")
			wDir1.On("IsDir").Return(true)
			wDir1.On("ModTime").Return(now.Add(time.Second * -100))

			wDir2 := new(FileInfoMock)
			wDir2.On("Name").Return("World two")
			wDir2.On("IsDir").Return(true)
			wDir2.On("ModTime").Return(now.Add(time.Second * -100))

			fsMock.On("ReadDir", "/home/world").Return([]os.FileInfo{wDir1, wDir2}, nil)

			f := data.Folder{Path: "/home/world"}
			world := f.AddWorld("World two")

			checkOneDir(w, &f)

			Convey("It should create a backup for a changed world and not unchanged ones", func() {
				fsMock.AssertExpectations(t)

				So(hasChangedFilesCallCount, ShouldEqual, 2)
				So(backupCreatedCallCount, ShouldEqual, 1)
				So(backedUpWorld.Id, ShouldNotBeEmpty)
				So(backedUpWorld.Id, ShouldNotEqual, world.Id)
				So(backedUpWorld.Name, ShouldEqual, "World one")

				Convey("and call checkPurgeBackup", func() {
					So(checkPurgeBackupCallCount, ShouldEqual, 1)
				})
			})
		})

	})

}

func TestWatcher_HasChangedFiles(t *testing.T) {
	Convey("Given a watcher and a world", t, func() {
		now := time.Now()

		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)

		world := data.World{
			Name:     "w1",
			FullPath: "/home/world/w1",
			Backups: []*data.Backup{
				{CreatedAt: now.Add(time.Second * -50)},
			},
		}

		f1 := new(FileInfoMock)
		f1.On("Name").Return("file1.txt")
		f1.On("IsDir").Return(false)
		f1.On("ModTime").Return(now.Add(time.Second * -100))

		f2 := new(FileInfoMock)
		f2.On("Name").Return("file2.txt")
		f2.On("IsDir").Return(false)

		Convey("When we fail to get world files", func() {
			fsMock.On("ReadDir", world.FullPath).Return([]os.FileInfo{}, errors.New("failed to read"))

			Convey("Then it should return false", func() {
				So(hasChangedFiles(log, fsMock, &world), ShouldBeFalse)
				fsMock.AssertExpectations(t)
			})
		})

		Convey("When we successfully get world files", func() {
			fsMock.On("ReadDir", world.FullPath).Return([]os.FileInfo{f1, f2}, nil)

			Convey("When there are no updated files since the last backup", func() {
				f2.On("ModTime").Return(now.Add(time.Second * -100))

				Convey("Then it should return false", func() {
					So(hasChangedFiles(log, fsMock, &world), ShouldBeFalse)
					fsMock.AssertExpectations(t)
				})
			})

			Convey("When there is an updated file", func() {
				f2.On("ModTime").Return(now)

				Convey("Then it should return true", func() {
					So(hasChangedFiles(log, fsMock, &world), ShouldBeTrue)
					fsMock.AssertExpectations(t)
				})
			})
		})
	})
}

func TestWatcher_CreateBackup(t *testing.T) {
	Convey("Given a watcher and a world to backup", t, func() {
		now := time.Unix(1495807405, 0)
		oldGetNow := getNow
		getNow = func() time.Time { return now }
		defer func() { getNow = oldGetNow }()

		config := conf.Config{
			BackupDir: "/back/up",
		}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)

		w := NewWatcher(log, &config, fsMock, dbMock)

		folder := data.Folder{Path: "/home/saves"}

		world := data.World{
			Id:       "WID01",
			Name:     "World One! For# Ever%Dude",
			FullPath: "/home/world/wee",
		}

		worldPath := path.Join(fmt.Sprintf(".%s", afero.FilePathSeparator), world.Name)

		Convey("When we fail to get current working directory", func() {
			fsMock.On("Getwd").Return("", errors.New("No!"))

			createBackup(w, log, &folder, &world)

			Convey("Then it should return", func() {
				fsMock.AssertExpectations(t)
			})
		})

		Convey("When the backup succeeds", func() {

			fsMock.On("Getwd").Return("/app/dir", nil)
			fsMock.On("Chdir", "/home/saves").Return(nil)
			fsMock.On("Chdir", "/app/dir").Return(nil)
			fsMock.On("Zip", worldPath, "/back/up/World_One_For_Ever_Dude-WID01-20170526T090325.zip").Return(nil)

			createBackup(w, log, &folder, &world)

			fsMock.AssertExpectations(t)

			Convey("Then it should add the backup to the world", func() {
				So(len(world.Backups), ShouldEqual, 1)
				So(world.Backups[0].Name, ShouldEqual, "World_One_For_Ever_Dude-WID01-20170526T090325.zip")
			})
		})

		Convey("When the backup fails", func() {
			fsMock.On("Getwd").Return("/app/dir", nil)
			fsMock.On("Chdir", "/home/saves").Return(nil)
			fsMock.On("Chdir", "/app/dir").Return(nil)
			fsMock.On("Zip", worldPath, "/back/up/World_One_For_Ever_Dude-WID01-20170526T090325.zip").Return(errors.New("Didn't work!"))

			createBackup(w, log, &folder, &world)

			fsMock.AssertExpectations(t)

			Convey("Then it should not add the backup to the world", func() {
				So(len(world.Backups), ShouldEqual, 0)
			})
		})
	})
}

func TestWatcher_CheckPurgeBackup(t *testing.T) {
	Convey("Given a valid watcher and world", t, func() {
		now := time.Now()
		oldGetNow := getNow
		getNow = func() time.Time { return now }
		defer func() { getNow = oldGetNow }()

		config := conf.Config{
			BackupDir:     "/back/up",
			CheckInterval: "5m",
		}
		log := logrus.WithField("test", "watcher")
		fsMock := new(IFileSystemMock)
		dbMock := new(IDbMock)

		w := NewWatcher(log, &config, fsMock, dbMock)

		world := data.World{
			Id:       "WID01",
			Name:     "w1",
			FullPath: "/home/world/wee",
		}

		Convey("When there are no backups", func() {
			Convey("It should do nothing", func() {
				checkPurgeBackup(w, log, &world)
			})
		})

		Convey("When there is only 1 backup", func() {
			world.Backups = []*data.Backup{
				{Id: "01", Name: "b1", CreatedAt: now},
			}

			Convey("It should do nothing", func() {
				checkPurgeBackup(w, log, &world)
			})
		})

		Convey("When there is more than 1 backup but they are older than our interval", func() {
			world.Backups = []*data.Backup{
				{Id: "01", Name: "b1", CreatedAt: now.Add(time.Minute * -20)},
				{Id: "02", Name: "b2", CreatedAt: now.Add(time.Minute * -15)},
				{Id: "03", Name: "b3", CreatedAt: now.Add(time.Minute * -10)},
			}

			Convey("It should do nothing", func() {
				checkPurgeBackup(w, log, &world)

				So(len(world.Backups), ShouldEqual, 3)
			})
		})

		Convey("When there is more than 1 backup and it is within our interval", func() {
			world.Backups = []*data.Backup{
				{Id: "01", Name: "b1", CreatedAt: now.Add(time.Minute * -20)},
				{Id: "02", Name: "b2", CreatedAt: now.Add(time.Minute * -15)},
				{Id: "03", Name: "b3", CreatedAt: now.Add(time.Minute * -10)},
				{Id: "04", Name: "b4", CreatedAt: now.Add(time.Minute * -5)},
				{Id: "05", Name: "b5", CreatedAt: now},
			}

			Convey("And the removal succeeds", func() {
				fsMock.On("Remove", "/back/up/b4").Return(nil)

				Convey("It should purge the backup", func() {
					checkPurgeBackup(w, log, &world)

					So(len(world.Backups), ShouldEqual, 4)
					fsMock.AssertExpectations(t)
				})
			})

			Convey("And the removal failes", func() {
				fsMock.On("Remove", "/back/up/b4").Return(errors.New("NO!"))

				Convey("It should not purge the backup", func() {
					checkPurgeBackup(w, log, &world)

					So(len(world.Backups), ShouldEqual, 5)
					fsMock.AssertExpectations(t)
				})
			})
		})
	})
}
