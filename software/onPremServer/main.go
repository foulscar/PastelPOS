package main

import (
	"log"
	"fmt"
	"net/http"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	const port int16 = 8080

	cliLogo := figure.NewColorFigure("PastelPOS", "larry3d", "purple", true)
	cliLogo.Print()
	
  log.Println("Loading Menu")
	menu, err := initMenu()
	if err != nil {
		log.Fatal(err)
	}
  
  log.Println("Loading Order Tracker System")
  ordersInSystem := initOrderTrackerSystem()

	log.Println("Starting WebSocket Server (Routed Through /socket.io)")
	wsServer := initWSServer(ordersInSystem);

	fs := http.FileServer(http.Dir("./static"))
	mux := http.NewServeMux()
	mux.Handle("/", fs)
	mux.Handle("/socket.io/", wsServer)
	mux.HandleFunc("/api/order", orderHandler(menu, ordersInSystem, wsServer))

	log.Println("Starting HTTP Server on Port: " + fmt.Sprint(port))
	log.Println("API Will Be Routed Through /api")
	log.Fatal(http.ListenAndServe(":" + fmt.Sprint(port), mux))
}
