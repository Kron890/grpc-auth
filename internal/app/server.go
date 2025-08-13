package app

import "github.com/labstack/echo/v4"

type Server struct {
	echo *echo.Echo
}

func NewServer() *Server {
	e := echo.New()
	return &Server{echo: e}
}
