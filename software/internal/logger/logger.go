package logger

import (
  "log"
  "fmt"
  "os"
)

func logService(service string, subService string, status string, messages ...string) {
  if len(messages) == 0 {
    return
  }

  log.Printf(
    "%s [%s] [%s] %s\n",
    service,
    subService,
    status,
    messages[0],
  )

  if len(messages) == 1 {
    return
  }

  for i := 1; i < len(messages); i++ {
    fmt.Println("\n", messages[i], "\n")
  }
}

func logServiceDEBUG(service string, subService string, messages ...string) {
  logService(service, subService, "DEBUG", messages...)
}

func logServiceINFO(service string, subService string, messages ...string) {
  logService(service, subService, "INFO", messages...)
}

func logServiceWARN(service string, subService string, messages ...string) {
  logService(service, subService, "WARN", messages...)
}

func logServiceERROR(service string, subService string, messages ...string) {
  logService(service, subService, "ERROR", messages...)
}

func logServiceFATAL(service string, subService string, messages ...string) {
  logService(service, subService, "FATAL", messages...)
  os.Exit(1)
}
