package main

import (
  "github.com/foulscar/PastelPOS/software/internal/logger"

  "github.com/common-nighthawk/go-figure"
)

func main() {
  cliLogo := figure.NewColorFigure("PastelPOS", "larry3d", "purple", true)
  cliLogo.Print()

  logServiceINFO("PastelPOS", "onPrem", "INFO", "Starting...")
}
