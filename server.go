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
  // take post body
  conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
  client := new(Client)
  if err != nil {
    log.Println(err.Error())
    return err
  }
  
  log.Println(conn.RemoteAddr().String(), "connected successfully")
  if err = c.Bind(client); err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err.Error())
  } 
  s.clients[conn.RemoteAddr().String()] = *client
  if (!(len(s.clients) > 1)) {
    return nil
  }
  fmt.Println(s.clients)
  carr := []Client{}
  for _, val := range s.clients {
    carr = append(carr, val)
  }
  roomId := RandSeq(5)
  carr[0].room = roomId
  carr[1].room = roomId
  room := Room{
  	c1:       carr[0],
  	c2:       carr[1],
  	messages: []Message{},
  }
  s.rooms[roomId] = room
  delete(s.clients, room.c1.conn.RemoteAddr().String())
  delete(s.clients, room.c2.conn.RemoteAddr().String())
  fmt.Println(s.rooms)
  fmt.Println(s.clients)
  // TODO: add matchmaking algorithm

  for {
    msg, _, err := wsutil.ReadClientData(conn)
    if err != nil {
      if err == io.EOF {
        log.Println(client.conn.RemoteAddr().String(), "disconnected")
        wsutil.WriteServerMessage(
          "User ",
        )
      }
    }
  }

}
