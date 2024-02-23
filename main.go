package main

import (
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
  e := echo.New()
  s := Server{
  	clients: map[*net.Conn]*Client{},
  	rooms:   map[string]Room{},
  }
  e.Use(middleware.Logger())
  e.GET("/", s.Serve)
  e.Logger.Fatal(e.Start(":8080"))
}
