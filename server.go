package main

import (
	"io"
	"log"
	"net"

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
  client.conn = conn
  for {
    msg, _, err := wsutil.ReadClientData(conn)
    if err != nil {
      log.Println(err)
      if err == io.EOF {
        if _, ok := s.clients[&conn]; ok {
          delete(s.clients, &conn)
          return nil
        }
        s.BroadcastMessage(client.room, Message{
        	Type: DSCNCTMSG,
        	From: *client,
        	Body: "",
        })
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
    } else if IsServerMessage(msg) {
      s.ParseServerMessage(msg, client, &conn)
      if len(s.clients) > 1 {
        s.Match()
      }
    } else {
      s.BroadcastMessage(
        client.room,
        Message{ Type: CLNTMSG, From: *client, Body: string(msg) },
      )
    }
  }
}

// TODO #2
// x 1. take connection, upgrade
// x 2. bind first message to client
// x 3. put clients in a separate room to wait for a match
// x 4. match top two after two clients connect
// x 5. inform client that match was made
// x 6. normal communication
// x 7. upon EOF, check if match was made, if yes remove from pool #1, or remove from pool #2 and inform other client
