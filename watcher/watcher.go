package watcher

import (
	"world-backup/conf"

	"fmt"
	"time"

	"world-backup/data"

	"os"

	"regexp"

	"github.com/Sirupsen/logrus"
	"errors"
)

var getNow = time.Now

type IArchiver interface {
	Make(zipPath string, filePaths []string) error
}

type IFileSystem interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
}

type Watcher struct {
	log    *logrus.Entry
	config *conf.Config
	fs     IFileSystem
	db     data.IDb
	zip    IArchiver
}

var NoWatchPathError = errors.New("No paths to watch")
var InvalidCheckInterval = errors.New("Invalid check interval")

func NewWatcher(log *logrus.Entry, config *conf.Config, fs IFileSystem, db data.IDb, zip IArchiver) *Watcher {
	w := Watcher{
		config: config,
		log:    log.WithField("component", "watcher"),
		fs:     fs,
		db:     db,
		zip:    zip,
	}

	return &w
}

func (w *Watcher) Start() error {
	if len(w.config.WatchDirs) == 0 {
		return NoWatchPathError
	}

	checkInterval, dErr := time.ParseDuration(w.config.CheckInterval)
	if dErr != nil {
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
		f := w.db.GetFolderByPath(w.config.WatchDirs[i])
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
			createBackup(w, worldLog, world)
		}
	}
}

var hasChangedFiles = func(log *logrus.Entry, fs IFileSystem, world *data.World) bool {
	// TODO: Need to add threshold so we are not continually backing up a world running.
	// TODO: Maybe if it is running, backup every 5 min (configurable) bot not every time
	// TODO: we poll
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

var createBackup = func(w *Watcher, log *logrus.Entry, world *data.World) {
	t := getNow()

	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	cleanWorldName := reg.ReplaceAllString(world.Name, "_")

	zipName := fmt.Sprintf("%s-%s-%s.zip", cleanWorldName, world.Id, t.Format("20060102T150405"))
	// TODO: Need to make OS env newline
	zipFullPath := fmt.Sprintf("%s/%s", w.config.BackupDir, zipName)

	log.Infof("Creating backup file %s", zipName)
	if err := w.zip.Make(zipFullPath, []string{world.FullPath}); err != nil {
		log.Errorf("Failed to  create zip: %s, %v", zipName, err)
		return
	}

	world.AddBackup(zipName)
}
