package watcher

import (
	"world-backup/data"

	"os"

	"time"

	"github.com/stretchr/testify/mock"
)

//region IDb
type IDbMock struct {
	mock.Mock
}

func (m *IDbMock) Save() error {
	args := m.Called()
	return args.Error(0)
}

func (m *IDbMock) Close() {
	m.Called()
}

func (m *IDbMock) AddFolder(path string) *data.Folder {
	args := m.Called(path)
	return args.Get(0).(*data.Folder)
}

func (m *IDbMock) Folders() []*data.Folder {
	args := m.Called()
	return args.Get(0).([]*data.Folder)
}

func (m *IDbMock) GetFolderByPath(path string) *data.Folder {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*data.Folder)
}

//endregion

//region IArchiver
type IArchiverMock struct {
	mock.Mock
}

func (m *IArchiverMock) Make(zipPath string, filePaths []string) error {
	args := m.Called(zipPath, filePaths)
	return args.Error(0)
}

//endregion

//region IFileSystem
type IFileSystemMock struct {
	mock.Mock
}

func (m *IFileSystemMock) ReadDir(dirname string) ([]os.FileInfo, error) {
	args := m.Called(dirname)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]os.FileInfo), args.Error(1)
}

func (m *IFileSystemMock) Remove(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

type FileInfoMock struct {
	mock.Mock
}

func (m *FileInfoMock) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *FileInfoMock) Size() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *FileInfoMock) Mode() os.FileMode {
	args := m.Called()
	return args.Get(0).(os.FileMode)
}

func (m *FileInfoMock) ModTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *FileInfoMock) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *FileInfoMock) Sys() interface{} {
	args := m.Called()
	return args.Get(0)
}

//endregion
