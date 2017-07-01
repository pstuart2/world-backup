package fs

import (
	"github.com/stretchr/testify/mock"
)

type IBackupFsMock struct {
	mock.Mock
}

func (m *IBackupFsMock) Chdir(dir string) error {
	args := m.Called(dir)
	return args.Error(0)
}

func (m *IBackupFsMock) Getwd() (dir string, err error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *IBackupFsMock) Zip(source, target string) error {
	args := m.Called(source, target)
	return args.Error(0)
}
