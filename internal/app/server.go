package app

import "github.com/labstack/echo/v4"

type Server struct {
	echo *echo.Echo
}

// NewServer cоздаем новый сервер через echo
func NewServer() *Server {
	e := echo.New()
	return &Server{echo: e}
}
