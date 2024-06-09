package onPrem

import (
  "log"

  socketio "github.com/googollee/go-socket.io"

  "github.com/foulscar/PastelPOS/software/internal/orders"
)

func NewWSServer(*orders.OrderManager) (*socketio.Server, error) {
  wsServer, err := socketio.NewServer(nil)
  if err != nil {
    return nil, err
  }
  
  wsServer.On("connection", func(so socketio.Socket) {
    so.On("join", func(room string) {
      so.Join(room)
      log.Println("WebSocket Server [ok]: A client has joined the room '" + room + "'")
    })

    so.On("requestInitPayload", func(room string) {
      log.Println("WebSocket Server [trying]: A client has requested initialization for the room '" + room + "'")
      EmitInitPayloadToClient(so, room)
    })
  })

  return wsServer, nil
}

func EmitInitPayloadToClient(so socketio.Socket, room string, orderManager *orders.OrderManager) {
  switch (room) {
  case "fohOrderTracker":
    EmitInitFOHOrderTracker(so, orderManager)
  }
}
