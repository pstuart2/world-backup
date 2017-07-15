package fs

import (
	"fmt"

	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/afero"
)

type IBackupFs interface {
	Getwd() (dir string, err error)
	Chdir(dir string) error
	Zip(source, target string) error
}

var CreateBackup = func(f IBackupFs, log *logrus.Entry, folderPath string, worldName string, backupDir string, backupName string) error {
	currentDir, dErr := f.Getwd()
	if dErr != nil {
		log.Errorf("Failed to get working dir: %v", dErr)
		return dErr
	}
	defer func() { f.Chdir(currentDir) }()
	f.Chdir(folderPath)

	zipFullPath := fmt.Sprintf("%s%s%s", backupDir, afero.FilePathSeparator, backupName)

	log.Infof("Creating backup file %s", zipFullPath)
	if err := f.Zip(worldName, zipFullPath); err != nil {
		log.Errorf("Failed to  create zip: %s, %v", zipFullPath, err)
		return err
	}

	return nil
}

func CleanName(name string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")
	return reg.ReplaceAllString(name, "_")
}
