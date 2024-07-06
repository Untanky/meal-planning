package database

import (
	"context"
	"database/sql"
	"errors"
	"meal-planning/domain"
	"strings"
	"time"
)

type mealDay struct {
	Date      string
	Breakfast string
	Lunch     string
	Dinner    string
	Snacks    *string
}

type sqlMealDayRepository struct {
	db *sql.DB
}

func NewSqlMealDayRepository(db *sql.DB) domain.MealDayRepository {
	return &sqlMealDayRepository{db}
}

func (s sqlMealDayRepository) FindByDate(ctx context.Context, date time.Time) (domain.MealDay, error) {
	row := s.db.QueryRowContext(ctx, `SELECT * FROM meals WHERE "date" = date(?) LIMIT 1`, date.Format("2006-01-02"))

	if row.Err() != nil {
		return domain.MealDay{}, row.Err()
	}

	day := new(mealDay)
	err := row.Scan(&day.Date, &day.Breakfast, &day.Lunch, &day.Dinner, &day.Snacks)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.MealDay{}, domain.MealNotFound
	} else if err != nil {
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

func (s sqlMealDayRepository) Create(ctx context.Context, mealDay domain.MealDay) (domain.MealDay, error) {
	_, err := s.db.ExecContext(ctx, `INSERT INTO meals VALUES (?, ?, ?, ?, ?)`, mealDay.Date.Format("2006-01-02"), mealDay.Breakfast, mealDay.Lunch, mealDay.Dinner, strings.Join(mealDay.Snacks, ","))

	if err != nil {
		return domain.MealDay{}, err
	}

	return mealDay, nil
}

func (s sqlMealDayRepository) Update(ctx context.Context, mealDay domain.MealDay) (domain.MealDay, error) {
	_, err := s.db.ExecContext(ctx, `UPDATE meals SET breakfast = ?, lunch = ?, dinner = ?, snacks = ? WHERE date = date(?)`, mealDay.Breakfast, mealDay.Lunch, mealDay.Dinner, strings.Join(mealDay.Snacks, ","), mealDay.Date.Format("2006-01-02"))

	if err != nil {
		return domain.MealDay{}, err
	}

	return mealDay, nil
}
