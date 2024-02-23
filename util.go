package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net"
	"strings"
	"time"
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

func (s *Server) match() {
  for k, v := range s.clients {
    for k2, v2 := range s.clients {
      if k == k2 {
        continue
      }
      distance := CalculateDistance(v.Location, v2.Location)
      if distance <= float64(v.Range) && distance <= float64(v2.Range) {
        // match dem
      }
    }
  }
}
