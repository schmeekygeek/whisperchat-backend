package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
  e := echo.New()
  s := Server{
  	clients: map[string]Client{},
  	rooms:   map[string]Room{},
  }
  e.POST("/", s.Serve)
  e.Logger.Fatal(e.Start(":8080"))
  log.Println("Listening on port 8080")
}
