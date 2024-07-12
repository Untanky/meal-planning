package domain

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

const (
	CaloriesPerKilogramBodyFat = 7700
)

var NutritionNotFound = errors.New("nutrition not found")

type Nutrition struct {
	Date     time.Time
	Calories *int
	Weight   *int
}

type AverageNutrition struct {
	Calories int
	Weight   int
}

type TotalDailyEnergyExpenditure struct {
	Start                       time.Time
	End                         time.Time
	AverageCalories             int
	PeriodWeightDifference      int
	TotalDailyEnergyExpenditure int
}

type NutritionRepository interface {
	FindByDate(ctx context.Context, date time.Time) (Nutrition, error)
	FindByDateRange(ctx context.Context, start, end time.Time) ([]Nutrition, error)
	FindAverageNutrition(ctx context.Context, start, end time.Time) (AverageNutrition, error)
	Create(ctx context.Context, n Nutrition) (Nutrition, error)
	Update(ctx context.Context, n Nutrition) (Nutrition, error)
	Delete(ctx context.Context, n Nutrition) error
}

type NutritionService struct {
	repository NutritionRepository
}

func NewNutritionService(repository NutritionRepository) *NutritionService {
	return &NutritionService{repository: repository}
}

func (service *NutritionService) FindByDateRange(ctx context.Context, start, end time.Time) ([]Nutrition, error) {
	slog.Info("Finding nutrition by date range", slog.String("start", start.Format("2006-01-02")), slog.String("end", end.Format("2006-01-02")))

	numberOfDays := int64(end.Sub(start).Hours() / 24)
	nutritionList := make([]Nutrition, 0, numberOfDays)

	dbNutritionList, err := service.repository.FindByDateRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		nutrition := Nutrition{
			Date: day,
		}

		for _, dbNutrition := range dbNutritionList {
			if dbNutrition.Date.Format("2006-01-02") == day.Format("2006-01-02") {
				nutrition = dbNutrition
			}
		}

		nutritionList = append(nutritionList, nutrition)
	}

	return nutritionList, nil
}

func (service *NutritionService) FindByDate(ctx context.Context, date time.Time) (Nutrition, error) {
	slog.Info("Finding nutrition by date", slog.String("date", date.Format("2006-01-02")))

	meal, err := service.repository.FindByDate(ctx, date)
	if errors.Is(MealNotFound, err) {
		return Nutrition{
			Date: date,
		}, nil
	}

	if err != nil {
		return Nutrition{}, err
	}

	return meal, nil
}

func (service *NutritionService) Upsert(ctx context.Context, nutrition Nutrition) (Nutrition, error) {
	slog.Info("Upserting dbNutrition", slog.String("date", nutrition.Date.Format("2006-01-02")))

	dbNutrition, err := service.repository.FindByDate(ctx, nutrition.Date)
	if err != nil && !errors.Is(err, NutritionNotFound) {
		return Nutrition{}, err
	}

	if errors.Is(err, NutritionNotFound) {
		slog.Debug("Nutrition does not exist", slog.String("date", nutrition.Date.Format("2006-01-02")))
		slog.Info("Creating Nutrition", slog.String("date", nutrition.Date.Format("2006-01-02")))

		dbNutrition, err = service.repository.Create(ctx, nutrition)
	} else {
		slog.Info("Updating Nutrition", slog.String("date", nutrition.Date.Format("2006-01-02")))

		dbNutrition, err = service.repository.Update(ctx, nutrition)
	}

	return dbNutrition, err
}

func (service *NutritionService) Delete(ctx context.Context, n Nutrition) error {
	panic("implement me")
}

func (service *NutritionService) CalculateTotalDailyEnergyExpenditure(ctx context.Context, start, end time.Time) (TotalDailyEnergyExpenditure, error) {
	previousPeriodStart := end.Add(-1 * 24 * time.Hour)
	previousPeriodEnd := previousPeriodStart.Add(-end.Sub(start))

	slog.Info("Calculating total daily energy expenditure")

	slog.Info("Finding average nutrition for previous period", slog.Time("start", previousPeriodStart), slog.Time("end", previousPeriodEnd))
	previousAverage, err := service.repository.FindAverageNutrition(ctx, previousPeriodStart, previousPeriodEnd)
	if err != nil {
		return TotalDailyEnergyExpenditure{}, err
	}

	slog.Info("Finding average nutrition for current period", slog.Time("start", start), slog.Time("end", end))
	currentAverage, err := service.repository.FindAverageNutrition(ctx, start, end)
	if err != nil {
		return TotalDailyEnergyExpenditure{}, err
	}

	weightDifference := currentAverage.Weight - previousAverage.Weight
	averageCalorieDifference := weightDifference * CaloriesPerKilogramBodyFat

	slog.Debug("Calculated average calorie difference", slog.Int("averageCalorieDifference", averageCalorieDifference), slog.Int("weightDifference", weightDifference))

	totalDailyEnergyExpenditure := currentAverage.Calories + averageCalorieDifference

	slog.Debug("Calculated total daily energy expenditure", slog.Int("totalDailyEnergyExpenditure", totalDailyEnergyExpenditure))

	return TotalDailyEnergyExpenditure{
		Start:                       start,
		End:                         end,
		AverageCalories:             currentAverage.Calories,
		PeriodWeightDifference:      weightDifference,
		TotalDailyEnergyExpenditure: totalDailyEnergyExpenditure,
	}, nil
}
