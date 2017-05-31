package data

import (
	"os"

	"github.com/stretchr/testify/mock"
)

//region IFileSystem
type IDbFileSystemMock struct {
	mock.Mock
}

func (m *IDbFileSystemMock) Exists(path string) (bool, error) {
	args := m.Called(path)
	return args.Bool(0), args.Error(1)
}

func (m *IDbFileSystemMock) ReadFile(filename string) ([]byte, error) {
	args := m.Called(filename)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]byte), args.Error(1)
}

func (m *IDbFileSystemMock) WriteFile(filename string, data []byte, perm os.FileMode) error {
	args := m.Called(filename, data, perm)
	return args.Error(0)
}

//endregion
