package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/echo/v4"
)

type Server struct {
  clients     map[string]Client
  rooms       map[string]Room
}

func (s *Server) Serve(c echo.Context) error {

  client := new(Client)
  fmt.Println("")
  if err := c.Bind(client); err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err.Error())
  }

  conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
  if err != nil {
    log.Println(err.Error())
    return err
  }
  
  log.Println(conn.RemoteAddr().String(), "connected successfully")
  if err = c.Bind(client); err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err.Error())
  } 
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
      fmt.Println("hi")
      parseServerMessage(string(msg))
    } else {
      sendClientMessage(string(msg), client.room)
    }
  }
}

func isServerMessage(msg string) bool {
  return len(msg) > 4 && msg[0:4] == "msg:"
}

func sendClientMessage(msg string, room string) {

}

func parseServerMessage(msg string) {
  message := msg[4:len(msg)]
  fmt.Println("server message", message, "received")
}

// TODO: matchmaking algorithm (match top two)
// TODO: bind client body after successfully upgrading
// TODO: remove connection from pool, inform other client and delete room
