package main

import "net"

type Location struct {
  Lat float64   `json:"lat"`
  Long float64   `json:"long"`
}

type Room struct {
  c1      Client
  c2      Client
  messages  []Message
}

type Message struct {
  From    Client `json:"from"`
  Body    string `json:"body"`
}

type Client struct {
  Username  string    `json:"username"`
  Location  Location  `json:"location"`
  Range     int       `json:"range"` // in kms
  conn      net.Conn
  room      string
  // isMatched bool
}

const (
  MATCHEDMSG = "msg:MATCHED"
  DISCONNECTED = "msg:DISCONNECTED %s"
)
