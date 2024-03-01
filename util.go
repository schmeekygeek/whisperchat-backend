package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/gobwas/ws/wsutil"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandSeq(n int) string {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  b := make([]rune, n)
  for i := range b {
    b[i] = letters[r.Intn(len(letters))]
  }
  return string(b)
}

func CalculateDistance(point1 Location, point2 Location) float64 {
  const R = 6371 //Radius of the earth in km
  dLat := deg2rad(float64(point1.Lat) - float64(point2.Lat))
  dLong := deg2rad(float64(point1.Long) - float64(point2.Long))
  a := math.Sin(dLat / 2) * math.Sin(dLat / 2) + math.Cos(
    deg2rad(float64(point1.Lat)) * math.Cos(deg2rad(float64(point2.Lat))),
  ) * math.Sin(dLong / 2) * math.Sin(dLong / 2)
  c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1 - a))
  d := R * c;
  return d
}

func deg2rad(deg float64) float64 {
  return deg * (math.Pi / 180)
}

func (s *Server) parseServerMessage(msg []byte, cl *Client, conn *net.Conn) {
  message := msg[4:len(string(msg))]
  if strings.HasPrefix(string(message), "BIND") {
    bodyToParse := message[5:]
    log.Println("Got body to parse:", string(bodyToParse))
    err := json.Unmarshal(bodyToParse, cl)
    if err != nil {
      log.Println(err.Error())
    }
    s.clients[conn] = cl
  }
}

func (s *Server) match() error {
  for k, v := range s.clients {
    for k2, v2 := range s.clients {
      if k == k2 {
        continue
      }
      distance := CalculateDistance(v.Location, v2.Location)
      if distance <= float64(v.Range) && distance <= float64(v2.Range) {
        // match dem
        room := new(Room)
        roomId := RandSeq(5)
        c1 := s.clients[k]
        c1.conn = *k
        c1.room = roomId
        room.c1 = *c1

        c2 := s.clients[k2]
        c2.conn = *k2
        c2.room = roomId
        room.c2 = *c2
        s.rooms[roomId] = *room
        // manually send the matched message to both clients
        s.sendClientDetails(*c1, *c2, MATCHEDMSG)
        s.sendClientDetails(*c2, *c1, MATCHEDMSG)
        delete(s.clients, k)
        delete(s.clients, k2)
        return nil
      }
    }
  }
  return nil
}

func (s *Server) broadcastMessage(roomId string, msg Message) {
  if room, ok := s.rooms[roomId]; ok {
    room.messages = append(room.messages, msg)
    s.rooms[roomId] = room
    jsn, err := json.Marshal(msg)
    if err != nil {
      log.Println(err.Error())
    }
    wsutil.WriteServerMessage(
      room.c1.conn,
      1,
      jsn,
    )
    wsutil.WriteServerMessage(
      room.c2.conn,
      1,
      jsn,
    )
  }
}

func (s *Server) sendClientDetails(to, of Client, msgType MessageType) {
  ofBodyJson, err := json.Marshal(&of)
  if err != nil {
    log.Println(err.Error())
  }
  msg := Message{
  	Type: msgType,
  	From: Client{},
  	Body: string(ofBodyJson),
  }
  msgJson, err := json.Marshal(&msg)
  if err != nil {
    log.Println(err.Error())
  }
  wsutil.WriteServerMessage(
    to.conn,
    1,
    msgJson,
  )
}

func isServerMessage(msg Message) bool {
  return msg.Type != CLNTMSG
}
