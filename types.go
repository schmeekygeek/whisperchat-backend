package main

import "net"

type Location struct {
  lat float64
  long float64
}

type Room struct {
  c1      Client
  c2      Client
  messages  []Message
}

type Message struct {
  from    Client
  body    string
}

type Client struct {
  room      string
  conn      net.Conn
  location  Location
  isMatched bool
}
