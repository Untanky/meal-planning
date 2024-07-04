package domain

import (
	"context"
	"time"
)

type MealDay struct {
	Date      time.Time
	Breakfast string
	Lunch     string
	Dinner    string
	Snacks    []string
}

type MealDayRepository interface {
	FindByDate(ctx context.Context, date time.Time) (MealDay, error)
	FindByDateRange(ctx context.Context, start, end time.Time) ([]MealDay, error)
	Create(ctx context.Context, mealDay MealDay) (MealDay, error)
	Update(ctx context.Context, mealDay MealDay) (MealDay, error)
}
