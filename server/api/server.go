package api

import (
	"net/http"

	"github.com/labstack/echo"
)

type IServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Start(string) error
	Use(middleware ...echo.MiddlewareFunc)
	Group(prefix string, m ...echo.MiddlewareFunc) (g IEchoGroup)
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	Routes() []echo.Route
	Static(prefix, root string)
}

type IEchoGroup interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
}

type EchoServer struct {
	e *echo.Echo
}

func (es EchoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	es.e.ServeHTTP(w, r)
}

func (es EchoServer) Static(prefix, root string) {
	es.e.Static(prefix, root)
}

func (es EchoServer) Start(address string) error {
	return es.e.Start(address)
}

func (es EchoServer) Use(middleware ...echo.MiddlewareFunc) {
	es.e.Use(middleware...)
}

func (es EchoServer) Group(prefix string, m ...echo.MiddlewareFunc) IEchoGroup {
	return es.e.Group(prefix, m...)
}

func (es EchoServer) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	es.e.GET(path, h, m...)
}

func (es EchoServer) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	es.e.POST(path, h, m...)
}

func (es EchoServer) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	es.e.POST(path, h, m...)
}

func (es EchoServer) Routes() []echo.Route {
	return es.e.Routes()
}
