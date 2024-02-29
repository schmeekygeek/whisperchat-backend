package main

import (
	"io"
	"log"
	"net"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/echo/v4"
)

type Server struct {
  clients     map[*net.Conn]*Client
  rooms       map[string]Room
}

func (s *Server) Serve(c echo.Context) error {

  conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
  if err != nil {
    log.Println(err.Error())
    return err
  }
  log.Println(conn.RemoteAddr().String(), "connected")
  client := new(Client)
  for {
    msg, _, err := wsutil.ReadClientData(conn)
    if err != nil {
      log.Println(err)
      if err == io.EOF {
        if v, ok := s.clients[&conn]; ok {
          log.Println("Unmatched")
          log.Println(v)
          delete(s.clients, &conn)
          return nil
        }
        s.broadcastMessage(client.room, DISCONNECTED)
        return nil
      }
    }
    if string(msg) == "hi" {
      // debug
      log.Println("\nRooms:", s.rooms)
      log.Println("___________")
      log.Println("Clients:")
      for k, v := range s.clients {
        log.Println(k, v)
      }
      log.Println("___________")
    } else if isServerMessage(string(msg)) {
      s.parseServerMessage(msg, client, &conn)
      if len(s.clients) > 1 {
        s.match()
      }
    } else {
      s.sendClientMessage(client, string(msg))
    }
  }
}

func Test(c echo.Context) error {
  conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
  log.Println(conn.RemoteAddr().String(), "connected")
  if err != nil {
    log.Println(err.Error())
  }

  for {
    msg, _, err := wsutil.ReadClientData(conn)
    if string(msg) == "hi" {
      if err != nil {
        log.Println(err.Error())
        return err
      }
    }
  }
}

func isServerMessage(msg string) bool {
  return strings.HasPrefix(msg, "msg:")
}

// TODO #2
// x 1. take connection, upgrade
// x 2. bind first message to client
// x 3. put clients in a separate room to wait for a match
// x 4. match top two after two clients connect
// x 5. inform client that match was made
// x 6. normal communication
// x 7. upon EOF, check if match was made, if yes remove from pool #1, or remove from pool #2 and inform other client
