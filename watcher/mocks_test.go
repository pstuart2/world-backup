package watcher

import (
	"world-backup/data"

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

func (m *IDbMock) Folders() []data.Folder {
	args := m.Called()
	return args.Get(0).([]data.Folder)
}

func (m *IDbMock) FolderByPath(path string) *data.Folder {
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
