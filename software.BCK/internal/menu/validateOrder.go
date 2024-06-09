package menu

import (
  "errors"

  "github.com/foulscar/PastelPOS/software/internal/orders"
)

func (menu Menu) validateOrder(order orders.Order) error {
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
