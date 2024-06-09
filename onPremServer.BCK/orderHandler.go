package main

import (
  "errors"
  "fmt"
  "log"
  "time"
	"net/http"
  "encoding/json"
  "math/rand"
  socketio "github.com/googollee/go-socket.io"
)

type MealItems struct {
  Primary string `json:"primary"`
  Secondary string `json:"secondary"`
  Drink string `json:"drink"`
}

type Meal struct {
	MealItems MealItems `json:"mealItems"`
	Count int8 `json:"count"`
}

type SingleItem struct {
  Type string `json:"type"`
	Item string `json:"item"`
	Count int8 `json:"count"`
}

type OrderIN struct {
  Meals []Meal `json:"meals"`
  SingleItems []SingleItem `json:"singleItems"`
}

type Order struct {
	Meals []Meal `json:"meals"`
	SingleItems []SingleItem `json:"singleItems"`
  Time int64 `json:"time"`
}

type RequestedOrderResponse struct {
  OrderID int16 `json:"orderID"`
}

type ItemNotFoundError struct {
  Item string
}

func (err *ItemNotFoundError) Error() string {
  return "Item '" + err.Item + "' not found in menu"
}

func (menu Menu) validateOrder(order OrderIN) error {
  if len(order.Meals) == 0 && len(order.SingleItems) == 0 {
    return errors.New("This order is empty")
  }

  for _, meal := range order.Meals {
    if _, exists := menu.ItemsAvailable.Primary[meal.MealItems.Primary]; !exists {
      return errors.New("Meal item: '" + meal.MealItems.Primary + "' is not a primary item")
    }
    if _, exists := menu.ItemsAvailable.Secondary[meal.MealItems.Secondary]; !exists {
      return errors.New("Meal item: '" + meal.MealItems.Secondary + "' is not a secondary item")
    }
    if _, exists := menu.ItemsAvailable.Drinks[meal.MealItems.Drink]; !exists {
      return errors.New("Meal item: '" + meal.MealItems.Drink + "' is not a drink item")
    }
  }
  
  for _, singleItem := range order.SingleItems {
    var exists bool

    switch singleItem.Type {
    case "primary":
      _, exists = menu.ItemsAvailable.Primary[singleItem.Item]
    case "secondary":
      _, exists = menu.ItemsAvailable.Secondary[singleItem.Item]
    case "drink":
      _, exists = menu.ItemsAvailable.Drinks[singleItem.Item]
    default:
      return errors.New(singleItem.Type + " is not a valid type")
    }

    if !exists {
      return errors.New("Single item: '" + singleItem.Item + "' is not a " + singleItem.Type + " item")
    }
  }

  return nil
}

func processOrder(orderIN OrderIN, ordersInSystem *OrdersInSystem, wsServer *socketio.Server) (orderID int16, err error) {
  if len(ordersInSystem.InProgress) >= 100 {
    return 0, errors.New("Too many orders")
  }

  var orderIDToGive int16
  for {
    tempID := int16(rand.Intn(9000) + 1000)

    if _, exists := ordersInSystem.InProgress[tempID]; !exists {
      orderIDToGive = tempID
      break
    } 
  }
  
  order := Order{
    Meals: orderIN.Meals,
    SingleItems: orderIN.SingleItems,
    Time: time.Now().Unix(),
  }

  orderPartialPrimary := OrderPartial{
    Items: make(map[string]int16),
    Time: order.Time,
  }
  orderPartialSecondary := OrderPartial{
    Items: make(map[string]int16),
    Time: order.Time,
  }

  for _, meal := range order.Meals {
    orderPartialPrimary.Items[meal.MealItems.Primary] += int16(meal.Count)
    orderPartialSecondary.Items[meal.MealItems.Secondary] += int16(meal.Count)
  }

  for _, singleItem := range order.SingleItems {
    switch singleItem.Type {
    case "primary":
      orderPartialPrimary.Items[singleItem.Item] += int16(singleItem.Count)
    case "secondary":
      orderPartialSecondary.Items[singleItem.Item] += int16(singleItem.Count)
    }
  }
 
  ordersInSystem.InProgress[orderIDToGive] = order
  if len(orderPartialPrimary.Items) > 0 {
    ordersInSystem.InProgressPrimary[orderIDToGive] = orderPartialPrimary
  }
  if len(orderPartialSecondary.Items) > 0 {
    ordersInSystem.InProgressSecondary[orderIDToGive] = orderPartialSecondary
  }
  
  log.Println("API: Order System has ingested an order")
  fmt.Println(ordersInSystem.InProgress[orderIDToGive])
  
  err = updateWSClientsNewOrder(
    wsServer,
    orderIDToGive,
    order,
    orderPartialPrimary,
    orderPartialSecondary,
  )
  if err != nil {
    delete(ordersInSystem.InProgress, orderIDToGive)
    return 0, err
  }

  return orderIDToGive, nil
}

func orderHandler(menu Menu, ordersInSystem *OrdersInSystem, wsServer *socketio.Server) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    var orderIN OrderIN

    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    err := decoder.Decode(&orderIN)
    if err != nil {
      http.Error(w, "Failed to parse request body", http.StatusBadRequest)
      return
    }

    err = menu.validateOrder(orderIN)
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

    orderID, err := processOrder(orderIN, ordersInSystem, wsServer)
    if err != nil {
      http.Error(w, err.Error(), http.StatusServiceUnavailable)
      return
    }
    
    requestedOrderResponse := RequestedOrderResponse{
      OrderID: orderID,
    }

    jsonResponseData, err := json.Marshal(requestedOrderResponse)
    if err != nil {
      http.Error(w, "Failed to encode JSON Response", http.StatusInternalServerError)
      return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "%s", jsonResponseData)
  }
}
