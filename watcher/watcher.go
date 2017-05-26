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
	w.check()

	go func() {
		shouldStop := false
		for !shouldStop {
			select {
			//case shouldStop = <-ci.stop:
			//	log.Infof("Quit message: %v", shouldStop)
			//	break
			case <-time.After(time.Second * 10):
				w.log.Info("Checking!")
				w.check()
			}
		}
	}()
}

func (w *Watcher) check() {
	dirs, err := w.fs.ReadDir(w.config.WatchDir)
	if err != nil {
		w.log.Error(err)
		return
	}

	t := getNow()
	w.log.Infof("Time: %d", t.Unix())

	for k, v := range dirs {
		w.log.Infof("%d - %s (isDir: %t) %d", k, v.Name(), v.IsDir(), v.ModTime().Unix())

		world := fmt.Sprintf("%s/%s", w.config.WatchDir, v.Name())
		zipName := fmt.Sprintf("%s/%s-%s.zip", w.config.BackupDir, v.Name(), t.Format("20060102T150405"))

		if ferr := w.zip.Make(zipName, []string{world}); ferr != nil {
			w.log.Errorf("Failed to  create zip: %s, %v", zipName, ferr)
		}
	}
}
