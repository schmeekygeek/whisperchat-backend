package main

import (
	"log"
	"net/http"
)

func main() {
  s := Server{
  	clients: map[string]Client{},
  	rooms:   map[string]Room{},
  }
  http.HandleFunc("/", s.Serve)
  log.Println("Listening on port 8080")
  http.ListenAndServe(":8080", nil)
}
