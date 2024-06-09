package main

import (
  "net/http"
  "github.com/gorilla/websocket"
)

type room struct {
  name string
  clients map[*client]bool
  join chan *client
  leave chan *client
  forward chan []byte
}

func newRoom(name string) *room {
  return &room{
    name: name,
    clients: make(map[*client]bool),
    join: make(chan *client),
    leave: make(chan *client),
    forward: make(chan []byte),
  }
}

func (r *room) run() {
  logService("WebSocket Server", r.name, "INFO", "Starting", nil)
  for {
    select {
    case client := <-r.join:
      r.clients[client] = true
    case client := <-r.leave:
      delete(r.clients, client)
      close(client.send)
    case msg := <-r.forward:
      for client := range r.clients {
        client.send <- msg
      }
    }
  }
}

const (
  socketBufferSize = 1024
  messageBufferSize = 256
)

var upgrader = websocket.Upgrader{
  ReadBufferSize: socketBufferSize,
  WriteBufferSize: socketBufferSize,
}

func (room *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  upgrader.CheckOrigin = func(r *http.Request) bool { return true }
  socket, err := upgrader.Upgrade(w, req, nil)
  if err != nil {
    logService("WebSocket Server", room.name, "ERROR", err.Error(), nil) 
  }
  client := &client{
    socket: socket,
    room: room,
    send: make(chan []byte, messageBufferSize),
  }
  room.join <- client
  defer func() { room.leave <- client }()
  go client.write()
  client.read()
}
