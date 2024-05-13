package main

import (
	"os"
	"io"
	"encoding/json"
)

type Item struct {
	Name string `json:"name"`
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

type OrdersInSystem struct {
  InProgress map[int16]Order `json:"inProgress"`
  InProgressPrimary map[int16]OrderPartial `json:"inProgressPrimary"`
  InProgressSecondary map[int16]OrderPartial `json:"inProgressSecondary"`
  Ready map[int16]Order `json:"ready"`
}

func initMenu() (Menu, error) {
	jsonFile, err := os.Open("menu.json")
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var menu Menu
	json.Unmarshal([]byte(byteValue), &menu)
	return menu, err
}

func initOrderTrackerSystem() *OrdersInSystem {
  ordersInSystem := new(OrdersInSystem)
  ordersInSystem.InProgress = make(map[int16]Order)
  ordersInSystem.InProgressPrimary = make(map[int16]OrderPartial)
  ordersInSystem.InProgressSecondary = make(map[int16]OrderPartial)
  return ordersInSystem
}
