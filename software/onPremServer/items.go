package main

import (
	"os"
	"io"
	"encoding/json"
)

type Item struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Price float32 `json:"price"`
}

type ItemsAvailable struct {
	Primary []Item `json:"primary"`
	Secondary []Item `json:"secondary"`
	Drinks []Item `json:"drinks"`
}

type Menu struct {
	ItemsAvailable ItemsAvailable `json:"itemsAvailable"`
	MealDiscount int8 `json:"mealDiscount"`
}

func initMenu() (Menu, error) {
	jsonFile, err := os.Open("menu.json")
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var menu Menu
	json.Unmarshal([]byte(byteValue), &menu)
	return menu, err
}
