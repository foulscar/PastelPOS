package main

import (
  "github.com/gorilla/websocket"
)

type client struct {
  socket *websocket.Conn
  room *room
  send chan []byte
}

func (c *client) read() {
  defer c.socket.Close()
  for {
    _, msg, err := c.socket.ReadMessage()
    if err != nil {
      return
    }
    msgString := string(msg)
    logService("WebSocket Server", c.room.name, "DEBUG", "Received Message", &msgString)
  }
}

func (c *client) write() {
  defer c.socket.Close()
  for msg := range c.send {
    err := c.socket.WriteMessage(websocket.TextMessage, msg)
    if err != nil {
      return
    }
  }
}
