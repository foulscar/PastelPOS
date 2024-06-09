package main

import (
	"log"
  "strconv"
  "encoding/json"
  "errors"
	socketio "github.com/googollee/go-socket.io"
)

type InitFOHOrderTrackerPayload struct {
  InProgress []int16 `json:"inProgress"`
  Ready []int16 `json:"ready"`
}

type InitFOHOrderPayload struct {
  Orders map[int16]Order `json:"orders"`
}

type WSClientNewOrder struct {
  OrderID int16 `json:"orderID"`
  Order Order `json:"order"`
}

type WSClientNewOrderPartial struct {
  OrderID int16 `json:"orderID"`
  Order OrderPartial `json:"order"`
}

func initWSServer(ordersInSystem *OrdersInSystem) *socketio.Server {
	wsServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	wsServer.On("connection", func(so socketio.Socket) {
		so.On("join", func(room string) {
			so.Join(room)
			log.Println("WebSocket Server: Client has Joined Room '" + room + "'")
		})

    so.On("requestInit fohOrderTracker", func() {
      log.Println("WebSocket Server: Initializing Client for fohOrderTracker")
      initWSClientFOHOrderTracker(so, ordersInSystem)
    })

    so.On("requestInit fohOrderBagger", func() {
      log.Println("WebSocket Server: Initializing Client for fohOrderBagger")
      initWSClientFOHOrderBagger(so, ordersInSystem)
    })

    so.On("removeOrder fohOrderBagger", func(orderID string) {
      log.Println("WebSocket Server: fohOrderBagger Client Sent 'removeOrder' for '" + orderID + "'");
      orderIDINT64, err := strconv.ParseInt(orderID, 10, 16)
      if err != nil {
        log.Println("Error Converting String to Int64 when Ingesting removeOrder from fohOrderBagger")
        return
      }
      orderIDINT16 := int16(orderIDINT64)

      wsClientRemoveOrderFOHOrderBagger(wsServer, so, ordersInSystem, orderIDINT16)
    })
	})

	return wsServer
}

func initWSClientFOHOrderTracker(so socketio.Socket, ordersInSystem *OrdersInSystem) {
  var inProgress = make([]int16, 0)
  var ready = make([]int16, 0)

  for id, _ := range ordersInSystem.InProgress {
    inProgress = append(inProgress, id)  
  }
  for id, _ := range ordersInSystem.Ready {
    ready = append(ready, id)  
  }

  initFOHOrderTrackerPayload := InitFOHOrderTrackerPayload{
    inProgress,
    ready,
  }
  
  initFOHOrderTrackerPayloadJSON, err := json.Marshal(initFOHOrderTrackerPayload)
  if err != nil {
    log.Println("Error Creating Initial Payload for fohOrderTracker: " + err.Error())
  }

  err = so.Emit("init", string(initFOHOrderTrackerPayloadJSON))
  if err != nil {
    log.Println("WebSocket Server: Error Initilaizing Client: " + err.Error())
  }
}

func initWSClientFOHOrderBagger(so socketio.Socket, ordersInSystem *OrdersInSystem) {
  initFOHOrderBaggerPayload := InitFOHOrderPayload{
    Orders: ordersInSystem.InProgress,
  }

  initFOHOrderBaggerPayloadJSON, err := json.Marshal(initFOHOrderBaggerPayload)
  if err != nil {
    log.Println("Error Creating Initial Payload for fohOrderBagger: " + err.Error())
  }

  err = so.Emit("init", string(initFOHOrderBaggerPayloadJSON))
  if err != nil {
    log.Println("WebSocket Server: Error Initilaizing Client: " + err.Error())
  }
}

func wsClientRemoveOrderFOHOrderBagger(wsServer *socketio.Server, so socketio.Socket, ordersInSystem *OrdersInSystem, orderID int16) {
  order := ordersInSystem.InProgress[orderID]
  ordersInSystem.Ready[orderID] = order

  delete(ordersInSystem.InProgress, orderID)
  delete(ordersInSystem.InProgressPrimary, orderID)
  delete(ordersInSystem.InProgressSecondary, orderID)

  var orderIDString = strconv.Itoa(int(orderID))
  so.Emit("removeOrder", orderIDString)
  wsServer.BroadcastTo("bohPrimary", "removeOrder", orderIDString)
  wsServer.BroadcastTo("bohSecondary", "removeOrder", orderIDString)
  wsServer.BroadcastTo("fohOrderTracker", "progressOrder", orderIDString)
}

func updateWSClientsNewOrder(
  wsServer *socketio.Server,
  orderID int16,
  order Order,
  orderPartialPrimary OrderPartial,
  orderPartialSecondary OrderPartial,
) error {
  wsClientNewOrder := WSClientNewOrder{
    OrderID: orderID,
    Order: order,
  }
  orderJSON, err := json.Marshal(wsClientNewOrder)
  if err != nil {
    log.Println(err)
    return errors.New("Failed To Parse System Interpreted Order into JSON")
  }

  wsClientNewOrderPrimary := WSClientNewOrderPartial{
    OrderID: orderID,
    Order: orderPartialPrimary,
  }
  orderPartialPrimaryJSON, err := json.Marshal(wsClientNewOrderPrimary)
  if err != nil {
    log.Println(err)
    return errors.New("Failed To Parse System Interpreted OrderPartialPrimary into JSON")
  }
  
  wsClientNewOrderSecondary := WSClientNewOrderPartial{
    OrderID: orderID,
    Order: orderPartialSecondary,
  }
  orderPartialSecondaryJSON, err := json.Marshal(wsClientNewOrderSecondary)
  if err != nil {
    log.Println(err)
    return errors.New("Failed To Parse System Interpreted OrderPartialSecondary into JSON")
  }

  log.Println("WebSocket Server: Sending clients new order: " + strconv.Itoa(int(orderID)))

  wsServer.BroadcastTo("fohOrderTracker", "addOrder", orderID)
  wsServer.BroadcastTo("fohOrderBagger", "addOrder", string(orderJSON))
  wsServer.BroadcastTo("bohPrimary", "addOrder", string(orderPartialPrimaryJSON))
  wsServer.BroadcastTo("bohSecondary", "addOrder", string(orderPartialSecondaryJSON))

  return nil
}
