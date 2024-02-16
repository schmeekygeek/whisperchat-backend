package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/echo/v4"
)

type Server struct {
  clients     map[string]Client
  rooms       map[string]Room
}

func (s *Server) Serve(c echo.Context) error {

  conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
  if err != nil {
    log.Println(err.Error())
    return err
  }
  log.Println(conn.RemoteAddr().String(), "connected successfully")

  client := new(Client)
  for {
    msg, _, err := wsutil.ReadClientData(conn)
    if err != nil {
      if err == io.EOF {
        log.Println(conn.RemoteAddr().String(), "disconnected")
        return nil
      }
      log.Println(err.Error())
    }
    if isServerMessage(string(msg)) {
      parseServerMessage(string(msg))
    } else {
      sendClientMessage(string(msg), client.room)
    }
  }
}

func isServerMessage(msg string) bool {
  return strings.HasPrefix(msg, "msg:")
}

func sendClientMessage(msg string, room string) {

}

func parseServerMessage(msg string) {
  message := msg[4:len(msg)]
  if strings.HasPrefix(message, "BIND") {
    bodyToParse := message[5:]
    log.Println("Got body to parse:", bodyToParse)
  }
  log.Println("Received message", message)
}

// TODO: matchmaking algorithm (match top two)
// TODO: bind client body after successfully upgrading
// TODO: remove connection from pool, inform other client and delete room
