package domain

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

type MealDay struct {
	Date      time.Time
	Breakfast string
	Lunch     string
	Dinner    string
	Snacks    []string
}

var MealNotFound = errors.New("meal: not found")

type MealDayRepository interface {
	FindByDate(ctx context.Context, date time.Time) (MealDay, error)
	FindByDateRange(ctx context.Context, start, end time.Time) ([]MealDay, error)
	Create(ctx context.Context, mealDay MealDay) (MealDay, error)
	Update(ctx context.Context, mealDay MealDay) (MealDay, error)
}

type MealDayService struct {
	repository MealDayRepository
}

func NewMealDayService(repository MealDayRepository) *MealDayService {
	return &MealDayService{repository: repository}
}

func (service *MealDayService) FindByDateRange(ctx context.Context, start, end time.Time) ([]MealDay, error) {
	slog.Info("Finding meals by date range", slog.String("start", start.Format("2006-01-02")), slog.String("end", end.Format("2006-01-02")))

	numberOfDays := int64(end.Sub(start).Hours() / 24)
	meals := make([]MealDay, 0, numberOfDays)

	dbMeals, err := service.repository.FindByDateRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		meal := MealDay{
			Date: day,
		}

		for _, dbMeal := range dbMeals {
			if dbMeal.Date.Format("2006-01-02") == day.Format("2006-01-02") {
				meal = dbMeal
			}
		}

		meals = append(meals, meal)
	}

	return meals, nil
}

func (service *MealDayService) FindByDate(ctx context.Context, date time.Time) (MealDay, error) {
	slog.Info("Finding meals by date", slog.String("date", date.Format("2006-01-02")))

	meal, err := service.repository.FindByDate(ctx, date)
	if errors.Is(MealNotFound, err) {
		return MealDay{
			Date: date,
		}, nil
	}

	if err != nil {
		return MealDay{}, err
	}

	return meal, nil
}

func (service *MealDayService) Upsert(ctx context.Context, mealDay MealDay) (MealDay, error) {
	slog.Info("Upserting meal", slog.String("date", mealDay.Date.Format("2006-01-02")))

	meal, err := service.repository.FindByDate(ctx, mealDay.Date)
	if err != nil && !errors.Is(err, MealNotFound) {
		return MealDay{}, err
	}

	if errors.Is(err, MealNotFound) {
		slog.Debug("Meal does not exist", slog.String("date", mealDay.Date.Format("2006-01-02")))
		slog.Info("Creating meal", slog.String("date", mealDay.Date.Format("2006-01-02")))

		meal, err = service.repository.Create(ctx, mealDay)
	} else {
		slog.Info("Updating meal", slog.String("date", mealDay.Date.Format("2006-01-02")))

		meal, err = service.repository.Update(ctx, meal)
	}

	return meal, err
}
