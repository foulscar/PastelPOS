package main

import (
  "log"
  "net/http"
  "github.com/common-nighthawk/go-figure"
)

func setupRoutes(orderSystem *orderSystem) {
  fs := http.FileServer(http.Dir("./web/onPrem"))
  http.Handle("/", fs)
  http.HandleFunc("/ws/fohOrderTracker", orderSystem.fohOrderTrackerRoom.ServeHTTP)
}

func main() {
  cliLogo := figure.NewColorFigure("PastelPOS", "larry3d", "purple", true)
  cliLogo.Print()
  
  orderSystem := initOrderSystem()
  setupRoutes(orderSystem)
  logService("HTTP Server", "MAIN", "INFO", "Starting", nil)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
