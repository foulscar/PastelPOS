package orders

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

type OrderPartial struct {
  Items map[string]int16 `json:"items"`
  Time int64 `json:"time"`
}

