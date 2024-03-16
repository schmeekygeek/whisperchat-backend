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

type MessageType string

type Message struct {
  Type    MessageType `json:"type"`
  From    Client `json:"from"`
  Body    string `json:"body"`
}

type Client struct {
  Username  string    `json:"username"`
  Location  Location  `json:"location"`
  Range     int       `json:"range"` // in kms
  conn      net.Conn
  room      string
}

const (
  CNNCTMSG = "connected"
  BINDMSG = "bind"
  DSCNCTMSG = "disconnected"
  CLNTMSG = "client"
)
