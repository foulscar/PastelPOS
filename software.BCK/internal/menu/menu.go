package menu

type Item struct {
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
