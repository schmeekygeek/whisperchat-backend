package main

import (
	"net/http"
)

type Server struct {
  clients     []Client
  rooms       map[string]Room
}

func (s Server) Serve(w http.ResponseWriter, r http.Request) {

}
