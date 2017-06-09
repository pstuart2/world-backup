package api

import (
	"net/http"

	"world-backup/data"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

//region ApiDb Mock
type ApiDbMock struct {
	mock.Mock
}

func (m *ApiDbMock) Folders() []*data.Folder {
	args := m.Called()
	return args.Get(0).([]*data.Folder)
}

//endregion

//region Echo Mock
type EchoServerMock struct {
	mock.Mock

	e *echo.Echo
}

func (em *EchoServerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	em.Called(w, r)
}

func (em *EchoServerMock) Start(port string) error {
	args := em.Called(port)
	return args.Error(0)
}

func (em *EchoServerMock) Static(prefix, root string) {
	em.Called(prefix, root)
}

func (em *EchoServerMock) Use(middleware ...echo.MiddlewareFunc) {
	em.Called(middleware)
}

func (em *EchoServerMock) Group(prefix string, m ...echo.MiddlewareFunc) (g IEchoGroup) {
	args := em.Called(prefix, m)
	return args.Get(0).(IEchoGroup)
}

func (em *EchoServerMock) Routes() []echo.Route {
	args := em.Called()
	return args.Get(0).([]echo.Route)
}

func (em *EchoServerMock) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

func (em *EchoServerMock) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

//endregion

//region Echo Group Mock
type EchoGroupMock struct {
	mock.Mock
}

func (em *EchoGroupMock) Use(middleware ...echo.MiddlewareFunc) {
	em.Called(middleware)
}

func (em *EchoGroupMock) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

func (em *EchoGroupMock) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

func (em *EchoGroupMock) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

func (em *EchoGroupMock) PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

func (em *EchoGroupMock) PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

//endregion

//region Handler Mock
type HandlerMock struct {
	mock.Mock
}

func (em *HandlerMock) Handler(next echo.Context) error {
	args := em.Called(next)
	return args.Error(0)
}

//endregion
