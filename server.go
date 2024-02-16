package main

import (
	"encoding/json"
	"io"
	"log"
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
      parseServerMessage(msg, client, c)
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

func parseServerMessage(msg []byte, cl *Client, c echo.Context) {
  message := msg[4:len(string(msg))]
  if strings.HasPrefix(string(message), "BIND") {
    bodyToParse := message[5:]
    log.Println("Got body to parse:", bodyToParse)
    err := json.Unmarshal(bodyToParse, cl)
    if err != nil {
      log.Println(err.Error())
    }
    log.Println(cl.Username)
  }
}

// TODO: matchmaking algorithm (match top two)
// TODO: remove connection from pool, inform other client and delete room
