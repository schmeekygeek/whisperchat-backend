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
        // inform other client if matched
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
    }
  }
}

func Test(c echo.Context) error {
  conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
  log.Println(conn.RemoteAddr().String(), "connected")
  test := 0
  if err != nil {
    log.Println(err.Error())
  }

  for {
    msg, _, err := wsutil.ReadClientData(conn)
    if string(msg) == "hi" {
      test++
      log.Println(test)
    }
    if err != nil {
      log.Println(err.Error())
      return err
    }
  }
}

func isServerMessage(msg string) bool {
  return strings.HasPrefix(msg, "msg:")
}

func (s *Server) broadcastMessage(roomId string, message string) {
  // TODO: inform client about match made
}

// TODO: remove connection from pool, inform other client and delete room

// TODO #2
// 1. take connection, upgrade
// 2. bind first message to client
// 3. put clients in a separate room to wait for a match
// 4. match top two after two clients connect
// 5. inform client that match was made
// 6. normal communication
// 7. upon EOF, check if match was made, if yes remove from pool #1, or remove from pool #2 and inform other client
// inform client about match made and other client disconnected
