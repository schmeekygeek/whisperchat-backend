package main

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
  e := echo.New()
  s := Server{
  	clients: map[string]Client{},
  	rooms:   map[string]Room{},
  }
  e.Use(middleware.Logger())
  e.GET("/", s.Serve)
  e.Logger.Fatal(e.Start(":8080"))
}
