package main

import (
	"fmt"
	"net/http"
)

type Meal struct {
	Items []string `json:"items"`
	Count int8 `json:"count"`
}

type SingleItem struct {
	Item string `json:"item"`
	Count int8 `json:"count"`
}

type OrderIN struct {
	Meals []Meal `json:"meals"`
	SingleItems []SingleItem `json:"singleItems"`
}

var  orderHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "test")
})
