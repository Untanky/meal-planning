package database

import (
	"context"
	"database/sql"
	"meal-planning/domain"
	"time"
)

type mealDay struct {
	Date      string
	Breakfast string
	Lunch     string
	Dinner    string
	Snacks    *string
}

type MealDayRepository interface {
	FindByDate(ctx context.Context, date time.Time) (domain.MealDay, error)
	FindByDateRange(ctx context.Context, start, end time.Time) ([]domain.MealDay, error)
	Create(ctx context.Context, mealDay mealDay) (mealDay, error)
	Update(ctx context.Context, mealDay mealDay) (mealDay, error)
}

type sqlMealDayRepository struct {
	db *sql.DB
}

func NewSqlMealDayRepository(db *sql.DB) MealDayRepository {
	return &sqlMealDayRepository{db}
}

func (s sqlMealDayRepository) FindByDate(ctx context.Context, date time.Time) (domain.MealDay, error) {
	row := s.db.QueryRowContext(ctx, `SELECT * FROM meals WHERE "date" = date(?) LIMIT 1`, date.Format("2006-01-02"))

	if row.Err() != nil {
		return domain.MealDay{}, row.Err()
	}

	day := new(mealDay)
	err := row.Scan(&day.Date, &day.Breakfast, &day.Lunch, &day.Dinner, &day.Snacks)
	if err != nil {
		return domain.MealDay{}, err
	}

	parsedDate, err := time.Parse("2006-01-02", day.Date)
	if err != nil {
		return domain.MealDay{}, err
	}

	return domain.MealDay{
		Date:      parsedDate,
		Breakfast: day.Breakfast,
		Lunch:     day.Lunch,
		Dinner:    day.Dinner,
		//Snacks:    strings.Split(*day.Snacks, ","),
	}, nil
}

func (s sqlMealDayRepository) FindByDateRange(ctx context.Context, start, end time.Time) ([]domain.MealDay, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT * FROM meals WHERE date >= date(?) AND date <= date(?)",
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]domain.MealDay, 0)
	for rows.Next() {
		meal := mealDay{}

		err = rows.Scan(&meal.Date, &meal.Breakfast, &meal.Lunch, &meal.Dinner, &meal.Snacks)
		if err != nil {
			return nil, err
		}

		date, err := time.Parse("2006-01-02", meal.Date)
		if err != nil {
			return nil, err
		}

		list = append(list, domain.MealDay{
			Date:      date,
			Breakfast: meal.Breakfast,
			Lunch:     meal.Lunch,
			Dinner:    meal.Dinner,
			//Snacks:    strings.Split(*meal.Snacks, ","),
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s sqlMealDayRepository) Create(ctx context.Context, mealDay mealDay) (mealDay, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlMealDayRepository) Update(ctx context.Context, mealDay mealDay) (mealDay, error) {
	//TODO implement me
	panic("implement me")
}
