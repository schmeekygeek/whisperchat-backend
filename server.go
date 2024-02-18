package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/echo/v4"
)

type Server struct {
  clients     []Client
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
  client.conn = conn
  for {
    msg, _, err := wsutil.ReadClientData(conn)
    if err != nil {
      log.Println(err.Error())
      if err == io.EOF {
        log.Println("huzaaaah", client)
        log.Println(conn.RemoteAddr().String(), "disconnected")
        if client.isMatched {
          room := s.rooms[client.room]
          log.Println("here tit is #1", room.c1.conn)
          log.Println("here tit is #2", room.c2.conn)
          var otherClient Client
          if room.c1.conn == nil {
            otherClient = room.c2
          } else {
            otherClient = room.c1
          }
          log.Println(otherClient.conn)
          wsutil.WriteServerMessage(
            otherClient.conn,
            1,
            []byte(fmt.Sprint(client.Username, "disconnected")),
          )
          // delete(s.rooms, client.room)
        } else {
          log.Println("2.     here are the rooms, sire", s.rooms)
          if len(s.clients) == 1 {
            s.clients = []Client{}
          }
        }
      }
      return nil
    }

    if isServerMessage(string(msg)) {
      s.parseServerMessage(msg, client, c)
      if len(s.clients) == 2 {
        s.match()
        client.isMatched = true
        log.Println("is he matcheeeeed????", client.isMatched)
      }
    } else {
      log.Println(client)
      sendClientMessage(string(msg), client.room, *client)
    }
  }
}

func isServerMessage(msg string) bool {
  return strings.HasPrefix(msg, "msg:")
}
 
func sendClientMessage(msg string, room string, client Client) {

}

func (s *Server) parseServerMessage(msg []byte, cl *Client, c echo.Context) {
  message := msg[4:len(string(msg))]
  if strings.HasPrefix(string(message), "BIND") {
    bodyToParse := message[5:]
    log.Println("Got body to parse:", string(bodyToParse))
    err := json.Unmarshal(bodyToParse, cl)
    if err != nil {
      log.Println(err.Error())
    }
    log.Println(cl.Username)
    s.clients = append(s.clients, *cl)
  }
}

func (s *Server) match() {
  roomId := RandSeq(5)
  room := Room{
    c1:       s.clients[0],
    c2:       s.clients[1],
    messages: []Message{},
  }
  room.c1.room = roomId
  room.c1.isMatched = true
  room.c2.room = roomId
  room.c2.isMatched = true
  s.rooms[roomId] = room
  log.Println(
    "Matched",
    room.c1.conn.RemoteAddr().String(),
    "and",
    room.c2.conn.RemoteAddr().String(),
    )
  s.broadcastMessage(roomId, "msg:MATCHED")
  log.Println("here are the rooms, sire", s.rooms)
}

func (s *Server) broadcastMessage(roomId string, message string) {
  // TODO: inform client abotu match made
}

// TODO: remove connection from pool, inform other client and delete room
