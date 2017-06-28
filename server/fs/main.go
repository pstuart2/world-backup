package fs

import (
	"os"

	"github.com/spf13/afero"
)

type FileSystem struct {
	af afero.Afero
}

func NewFs(fs afero.Fs) *FileSystem {
	f := FileSystem{
		afero.Afero{Fs: fs},
	}

	return &f
}

func (f *FileSystem) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (f *FileSystem) Getwd() (dir string, err error) {
	return os.Getwd()
}

func (f *FileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return f.af.ReadDir(dirname)
}

func (f *FileSystem) Remove(name string) error {
	return f.af.Remove(name)
}

func (f *FileSystem) Exists(path string) (bool, error) {
	return f.af.Exists(path)
}

func (f *FileSystem) Rename(oldname, newname string) error {
	return f.af.Rename(oldname, newname)
}
