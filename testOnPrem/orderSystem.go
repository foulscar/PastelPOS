package main

import (
	"os"
	"io"
	"encoding/json"
)

type item struct {
	Name string `json:"name"`
  Icon string `json:"icon"`
	Price float32 `json:"price"`
}

type ItemsAvailable struct {
	Primary map[string]Item `json:"primary"`
	Secondary map[string]Item `json:"secondary"`
	Drinks map[string]Item `json:"drinks"`
}

type Menu struct {
	ItemsAvailable ItemsAvailable `json:"itemsAvailable"`
	MealDiscount int8 `json:"mealDiscount"`
}

type OrderPartial struct {
  Items map[string]int16 `json:"items"`
  Time int64 `json:"time"`
}

type orderSystem struct {
  fohOrderTrackerRoom *room
}

func initOrderSystem() *orderSystem {
  orderSystem := orderSystem{
    fohOrderTrackerRoom: newRoom("FOHOrderTracker"),
  }
  go orderSystem.fohOrderTrackerRoom.run()
  return &orderSystem
}
