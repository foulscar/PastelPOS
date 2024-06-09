package main

import (
  "log"
)

func logService(service string, subService string, status string, messageShort string, messageLong *string) {
  stdOutMessageLong := ""
  if messageLong != nil {
    stdOutMessageLong = "\n\n" + *messageLong
  }
  log.Printf(
    "%s [%s] [%s]: %s%s\n",
    service,
    subService,
    status,
    messageShort,
    stdOutMessageLong,
  ) 
}
