package watcher

import (
	"world-backup/conf"

	"fmt"
	"time"

	"world-backup/data"

	"os"

	"github.com/Sirupsen/logrus"
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

func NewWatcher(log *logrus.Entry, config *conf.Config, fs IFileSystem, db data.IDb, zip IArchiver) *Watcher {
	w := &Watcher{
		config: config,
		log:    log.WithField("component", "watcher"),
		fs:     fs,
		db:     db,
		zip:    zip,
	}

	return w
}

func (w *Watcher) Start() {
	if len(w.config.WatchDirs) == 0 {
		return
	}

	for i, d := range w.config.WatchDirs {
		w.log.Infof("Checking tracking for dir (%d) [%s]", i, d)

		f := w.db.FolderByPath(d)
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
	go watch(w, stopChannel)
}

var watch = func(w *Watcher, stop chan bool) {
	d, _ := time.ParseDuration(fmt.Sprintf("%ds", w.config.CheckIntervalSeconds))
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
	dirs, err := w.fs.ReadDir(w.config.WatchDirs[0])
	if err != nil {
		w.log.Error(err)
		return
	}

	t := getNow()
	w.log.Infof("Time: %d", t.Unix())

	for k, v := range dirs {
		w.log.Infof("%d - %s (isDir: %t) %d", k, v.Name(), v.IsDir(), v.ModTime().Unix())

		world := fmt.Sprintf("%s/%s", w.config.WatchDirs[0], v.Name())
		zipName := fmt.Sprintf("%s/%s-%s.zip", w.config.BackupDir, v.Name(), t.Format("20060102T150405"))

		if ferr := w.zip.Make(zipName, []string{world}); ferr != nil {
			w.log.Errorf("Failed to  create zip: %s, %v", zipName, ferr)
		}
	}
}
