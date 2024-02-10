package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Server struct {
  clients     map[string]Client
  rooms       map[string]Room
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {

  lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
  if err != nil {
    w.Write([]byte("Can't parse value for latitude"))
    // w.WriteHeader(400)
    // TODO: take distance preference from query
    return
  }
  long, err := strconv.ParseFloat(r.URL.Query().Get("long"), 64)
  if err != nil {
    w.Write([]byte("Can't parse value for longitude"))
    w.WriteHeader(400)
    return
  }
  location := Location{
  	lat:  lat,
  	long: long,
  }

  conn, _, _, err := ws.UpgradeHTTP(r, w)

  if err != nil {
    log.Println(err.Error())
    return
  }
  
  log.Println(conn.RemoteAddr().String(), "connected successfully")

  client := Client{
  	room:      "",
  	conn:      conn,
  	location:  location,
  	isMatched: true,
  }
  s.clients[conn.RemoteAddr().String()] = client
  if (!(len(s.clients) > 1)) {
    return
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
        
      }
    }
  }

}
