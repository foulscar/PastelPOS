package onPrem

import (
  "encoding/json"
  "log"

  socketio "github.com/googollee/go-socket.io"

  "github.com/foulscar/PastelPOS/software/internal/orders"
)

type InitFOHOrderTrackerPayload struct {
  InProgress []int16 `json:"inProgress"`
  Ready []int16 `json:"ready"`
}

type InitFOHOrderPayload struct {
  Orders map[int16]orders.Order `json:"orders"`
}

type WSClientNewOrder struct {
  OrderID int16 `json:"orderID"`
  Order orders.Order `json:"order"`
}

type WSClientNewOrderPartial struct {
  OrderID int16 `json:"orderID"`
  Order orders.OrderPartial `json:"order"`
}

func EmitInitFOHOrderTracker(so socketio.Socket, orderManager *orders.OrderManager) {
  var inProgress = make([]int16, 0)
  var ready = make([]int16, 0)

  for id, _ := range orderManager.OrdersInProgress {
    inProgress = append(inProgress, id)  
  }
  for id, _ := range orderManager.OrdersReady {
    ready = append(ready, id)  
  }

  initFOHOrderTrackerPayload := InitFOHOrderTrackerPayload{
    inProgress,
    ready,
  }
  
  initFOHOrderTrackerPayloadJSON, err := json.Marshal(initFOHOrderTrackerPayload)
  if err != nil {
    log.Println("WebSocket Server [fail]: Error Creating Initial Payload for fohOrderTracker: " + err.Error())
  }

  err = so.Emit("init", string(initFOHOrderTrackerPayloadJSON))
  if err != nil {
    log.Println("WebSocket Server [fail]: Error Initilaizing Client for fohOrderTracker: " + err.Error())
  }

  log.Println("WebSocket Server [ok]: Initialized client for fohOrderTracker")
}

func EmitInitFOHOrderBagger(so socketio.Socket, orderManager *orders.OrderManager) {
  initFOHOrderBaggerPayload := InitFOHOrderPayload{
    Orders: orderManager.OrdersInProgress,
  }

  initFOHOrderBaggerPayloadJSON, err := json.Marshal(initFOHOrderBaggerPayload)
  if err != nil {
    log.Println("WebSocket Server [fail]: Error Creating Initial Payload for fohOrderBagger: " + err.Error())
  }

  err = so.Emit("init", string(initFOHOrderBaggerPayloadJSON))
  if err != nil {
    log.Println("WebSocket Server [fail]: Error Initilaizing Client for fohOrderBagger: " + err.Error())
  }

  log.Println("WebSocket Server [ok]: Initialized client for fohOrderBagger")
}
