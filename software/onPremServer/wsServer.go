package main

import (
	"log"
	socketio "github.com/googollee/go-socket.io"
)

func initWSServer() *socketio.Server {
	wsServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	wsServer.On("connection", func(so socketio.Socket) {
		log.Println("WebSocket Server: Received a Connection")

		so.On("join", func(room string) {
			so.Join(room)
			log.Println("WebSocket Server: Client has Joined Room '" + room + "'")
		})

		so.On("test", func() { wsServer.BroadcastTo("fohOrderTracker", "addOrder", "123") })
		so.On("test2", func() { wsServer.BroadcastTo("fohOrderTracker", "progressOrder", "123") })
		so.On("test3", func() { wsServer.BroadcastTo("fohOrderTracker", "removeOrder", "123") })
	})

	return wsServer
}
