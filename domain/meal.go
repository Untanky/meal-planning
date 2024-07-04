package domain

import "time"

type MealDay struct {
	Date      time.Time
	Breakfast string
	Lunch     string
	Dinner    string
	Snacks    []string
}
