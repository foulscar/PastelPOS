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

		so.On("test", func() { wsServer.BroadcastTo("fohOrderTracker", "addOrder", "123") })
		so.On("test2", func() { wsServer.BroadcastTo("fohOrderTracker", "progressOrder", "123") })
		so.On("test3", func() { wsServer.BroadcastTo("fohOrderTracker", "removeOrder", "123") })
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
    log.Println("Error Creating Initial Payload for fohOrderTracker: ", err)
  }

  err = so.Emit("init", string(initFOHOrderTrackerPayloadJSON))
  if err != nil {
    log.Println("WebSocket Server: Error Initilaizing Client: " + err.Error())
  }
}

func updateWSClientsNewOrder(
  wsServer *socketio.Server,
  orderID int16,
  order Order,
  orderPartialPrimary OrderPartial,
  orderPartialSecondary OrderPartial,
) error {
  orderJSON, err := json.Marshal(order)
  if err != nil {
    log.Println(err)
    return errors.New("Failed To Parse System Interpreted Order into JSON")
  }

  orderPartialPrimaryJSON, err := json.Marshal(orderPartialPrimary)
  if err != nil {
    log.Println(err)
    return errors.New("Failed To Parse System Interpreted OrderPartialPrimary into JSON")
  }
  
  orderPartialSecondaryJSON, err := json.Marshal(orderPartialSecondary)
  if err != nil {
    log.Println(err)
    return errors.New("Failed To Parse System Interpreted OrderPartialSecondary into JSON")
  }

  log.Println("WebSocket Server: Sending clients new order: " + strconv.Itoa(int(orderID)))

  wsServer.BroadcastTo("fohOrderTracker", "addOrder", orderID)
  wsServer.BroadcastTo("fohOrderBagger", "addOrder", orderJSON)
  wsServer.BroadcastTo("bohPrimary", "addOrder", orderPartialPrimaryJSON)
  wsServer.BroadcastTo("bohSecondary", "addOrder", orderPartialSecondaryJSON)

  return nil
}
