package orders

import (
  "errors"
  "math/rand"
)

type OrderManager struct {
  OrdersInProgress map[int16]Order
  OrdersInProgressPrimary map[int16]OrderPartial
  OrdersInProgressSecondary map[int16]OrderPartial
  OrdersReady map[int16]Order

  OrderLimit int

  IngestedOrderHandler func(orderID int16, order Order) 
  IngestedOrderPrimaryHandler func(orderID int16, orderPartial OrderPartial)
  IngestedOrderSecondaryHandler func(orderID int16, orderPartial OrderPartial)
}

func NewOrderManager(
  orderLimit int,
  ingestedOrderHandler func(orderID int16, order Order),
  ingestedOrderPrimaryHandler func(orderID int16, orderPartial OrderPartial),
  ingestedOrderSecondaryHandler func(orderID int16, orderPartial OrderPartial),
) *OrderManager {
  orderManager := new(OrderManager)
  orderManager.OrdersInProgress = make(map[int16]Order)
  orderManager.OrdersInProgressPrimary = make(map[int16]OrderPartial)
  orderManager.OrdersInProgressSecondary = make(map[int16]OrderPartial)
  orderManager.OrdersReady = make(map[int16]Order)

  orderManager.OrderLimit = orderLimit

  orderManager.IngestedOrderHandler = ingestedOrderHandler
  orderManager.IngestedOrderPrimaryHandler = ingestedOrderPrimaryHandler
  orderManager.IngestedOrderSecondaryHandler = ingestedOrderSecondaryHandler
  
  return orderManager
}

func (orderManager OrderManager) ProcessOrder(order Order) (orderID int16, err error) {
  if len(orderManager.OrdersInProgress) >= orderManager.OrderLimit {
    return 0, errors.New("Too many orders")
  }

  var orderIDToGive int16
  for {
    tempID := int16(rand.Intn(9000) + 1000)
    _, inProgressExists := orderManager.OrdersInProgress[tempID]
    _, readyExists := orderManager.OrdersReady[tempID]

    if !inProgressExists && !readyExists {
      orderIDToGive = tempID
      break
    }
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
 
  orderManager.OrdersInProgress[orderIDToGive] = order

  if len(orderPartialPrimary.Items) > 0 {
    orderManager.OrdersInProgressPrimary[orderIDToGive] = orderPartialPrimary
    if orderManager.IngestedOrderPrimaryHandler != nil {
      go orderManager.IngestedOrderPrimaryHandler(orderIDToGive, orderPartialPrimary)
    }
  }

  if len(orderPartialSecondary.Items) > 0 {
    orderManager.OrdersInProgressSecondary[orderIDToGive] = orderPartialSecondary
    if orderManager.IngestedOrderSecondaryHandler != nil {
      go orderManager.IngestedOrderSecondaryHandler(orderIDToGive, orderPartialSecondary)
    }
  }
  
  if orderManager.IngestedOrderHandler != nil {
    go orderManager.IngestedOrderHandler(orderIDToGive, order)
  }

  return orderIDToGive, nil
}
