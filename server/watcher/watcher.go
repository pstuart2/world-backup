package watcher

import (
	"world-backup/server/conf"

	"fmt"
	"time"

	"world-backup/server/data"

	"os"

	"regexp"

	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/afero"
)

var getNow = time.Now

type IFileSystem interface {
	Chdir(dir string) error
	Getwd() (dir string, err error)
	ReadDir(dirname string) ([]os.FileInfo, error)
	Remove(name string) error
	Zip(source, target string) error
}

type IDb interface {
	Save() error
	Close()

	AddFolder(path string) *data.Folder
	GetFolderByPath(path string) *data.Folder
}

type Watcher struct {
	log    *logrus.Entry
	config *conf.Config
	fs     IFileSystem
	db     IDb
}

var NoWatchPathError = errors.New("No paths to watch")
var InvalidCheckInterval = errors.New("Invalid check interval")
var InvalidMinBackupAge = errors.New("Invalid min backup age")

func NewWatcher(log *logrus.Entry, config *conf.Config, fs IFileSystem, db IDb) *Watcher {
	w := Watcher{
		config: config,
		log:    log.WithField("component", "watcher"),
		fs:     fs,
		db:     db,
	}

	return &w
}

func (w *Watcher) Start() error {
	if len(w.config.WatchDirs) == 0 {
		return NoWatchPathError
	}

	checkInterval, ciError := time.ParseDuration(w.config.CheckInterval)
	if ciError != nil {
		return InvalidCheckInterval
	}

	for i, d := range w.config.WatchDirs {
		w.log.Infof("Checking tracking for dir (%d) [%s]", i, d)

		f := w.db.GetFolderByPath(d)
		if f == nil {
			w.log.Infof("Creating tracking for [%s}", d)
			f = w.db.AddFolder(d)
		}

		w.log.Infof("Watching: %s: %s", f.Id, f.Path)
	}

	w.db.Save()

	// Run our check right at startup
	check(w)

	stopChannel := make(chan bool)
	go watch(w, stopChannel, checkInterval)

	return nil
}

var watch = func(w *Watcher, stop chan bool, d time.Duration) {
	shouldStop := false
	for !shouldStop {
		select {
		case shouldStop = <-stop:
			w.log.Infof("Quit message: %v", shouldStop)
		case <-time.After(d):
			w.log.Info("Checking!")
			check(w)
		}
	}
}

var check = func(w *Watcher) {
	for i := range w.config.WatchDirs {
		path := w.config.WatchDirs[i]
		f := w.db.GetFolderByPath(path)
		checkOneDir(w, f)
		f.LastRun = getNow()
		w.db.Save()
	}
}

var checkOneDir = func(w *Watcher, f *data.Folder) {
	log := w.log.WithField("folder", f.Id)

	worldDirs, err := w.fs.ReadDir(f.Path)
	if err != nil {
		log.Error(err)
		return
	}

	for k, v := range worldDirs {
		log.Infof("%d - %s (isDir: %t) %d", k, v.Name(), v.IsDir(), v.ModTime().Unix())

		world := f.GetWorldByName(v.Name())
		if world == nil {
			world = f.AddWorld(v.Name())
		}

		worldLog := log.WithField("world", world.Id)
		if hasChangedFiles(worldLog, w.fs, world) {
			createBackup(w, worldLog, f, world)
			checkPurgeBackup(w, worldLog, world)
		}
	}
}

var hasChangedFiles = func(log *logrus.Entry, fs IFileSystem, world *data.World) bool {
	files, err := fs.ReadDir(world.FullPath)
	if err != nil {
		log.Errorf("Failed to check world [%s] for changes: %v", world.FullPath, err)
		return false
	}

	lastBackupTime := world.LastBackupTime()
	log.Infof("Last backup time: %d", lastBackupTime.Unix())

	for i := range files {
		file := files[i]
		if lastBackupTime.Before(file.ModTime()) {
			log.Infof("%s file was changed at %d", file.Name(), file.ModTime().Unix())
			return true
		}
	}

	return false
}

var createBackup = func(w *Watcher, log *logrus.Entry, f *data.Folder, world *data.World) {
	t := getNow()

	currentDir, dErr := w.fs.Getwd()
	if dErr != nil {
		log.Errorf("Failed to change workind dir: %v", dErr)
		return
	}
	defer func() { w.fs.Chdir(currentDir) }()
	w.fs.Chdir(f.Path)

	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	cleanWorldName := reg.ReplaceAllString(world.Name, "_")

	zipName := fmt.Sprintf("%s-%s-%s.zip", cleanWorldName, world.Id, t.Format("20060102T150405"))
	zipFullPath := fmt.Sprintf("%s%s%s", w.config.BackupDir, afero.FilePathSeparator, zipName)

	log.Infof("Creating backup file %s", zipName)
	if err := w.fs.Zip(world.Name, zipFullPath); err != nil {
		log.Errorf("Failed to  create zip: %s, %v", zipName, err)
		return
	}

	world.AddBackup(zipName)
}

var checkPurgeBackup = func(w *Watcher, log *logrus.Entry, world *data.World) {
	if len(world.Backups) < 2 {
		return
	}

	now := getNow()
	previousBackup := world.Backups[len(world.Backups)-2]
	checkInterval, _ := time.ParseDuration(w.config.CheckInterval)
	checkIntervalWithBuffer := checkInterval + (time.Second * 2)

	if previousBackup.CreatedAt.After(now.Add(-checkIntervalWithBuffer)) {
		zipName := fmt.Sprintf("%s%s%s", w.config.BackupDir, afero.FilePathSeparator, previousBackup.Name)
		log.Infof("Removing previous backup (%s) %s", previousBackup.Id, zipName)
		if err := w.fs.Remove(zipName); err != nil {
			log.Errorf("Failed to remove previous backup (%s), err: %v", zipName, err)
			return
		}

		world.RemoveBackup(previousBackup.Id)
	}
}
