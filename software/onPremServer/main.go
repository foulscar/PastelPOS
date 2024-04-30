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
	
	menu, err := initMenu()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(menu.ItemsAvailable.Primary[0])

	log.Println("Starting WebSocket Server (Routed Through /socket.io)")
	wsServer := initWSServer();

	fs := http.FileServer(http.Dir("./static"))
	mux := http.NewServeMux()
	mux.Handle("/", fs)
	mux.Handle("/socket.io/", wsServer)
	mux.Handle("/api/order", orderHandler)

	log.Println("Starting HTTP Server on Port: " + fmt.Sprint(port))
	log.Println("API Will Be Routed Through /api")
	log.Fatal(http.ListenAndServe(":" + fmt.Sprint(port), mux))
}
