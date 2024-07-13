package database

import (
	"context"
	"database/sql"
	"errors"
	"meal-planning/domain"
	"time"
)

type nutritionEntity struct {
	date     string
	calories sql.NullInt64
	weight   sql.NullInt64
}

type sqlNutritionRepository struct {
	db *sql.DB
}

func NewSqlNutritionRepository(db *sql.DB) domain.NutritionRepository {
	return &sqlNutritionRepository{
		db: db,
	}
}

func (s *sqlNutritionRepository) FindByDate(ctx context.Context, date time.Time) (domain.Nutrition, error) {
	row := s.db.QueryRowContext(ctx, `SELECT * FROM nutrition WHERE "date" = date(?) LIMIT 1`, date.Format("2006-01-02"))

	if row.Err() != nil {
		return domain.Nutrition{}, row.Err()
	}

	entity := new(nutritionEntity)
	err := row.Scan(&entity.date, &entity.calories, &entity.weight)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Nutrition{}, domain.NutritionNotFound
	} else if err != nil {
		return domain.Nutrition{}, err
	}

	parsedDate, err := time.Parse("2006-01-02", entity.date)
	if err != nil {
		return domain.Nutrition{}, err
	}

	var calories *int
	if entity.calories.Valid {
		temp := int(entity.calories.Int64)
		calories = &temp
	}

	var weight *int
	if entity.weight.Valid {
		temp := int(entity.weight.Int64)
		weight = &temp
	}

	return domain.Nutrition{
		Date:     parsedDate,
		Calories: calories,
		Weight:   weight,
	}, nil
}

func (s *sqlNutritionRepository) FindByDateRange(ctx context.Context, start, end time.Time) ([]domain.Nutrition, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT * FROM nutrition WHERE date >= date(?) AND date <= date(?)",
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]domain.Nutrition, 0)
	for rows.Next() {
		entity := nutritionEntity{}
		err = rows.Scan(&entity.date, &entity.calories, &entity.weight)
		if err != nil {
			return nil, err
		}

		date, err := time.Parse("2006-01-02", entity.date)
		if err != nil {
			return nil, err
		}

		var calories *int
		if entity.calories.Valid {
			temp := int(entity.calories.Int64)
			calories = &temp
		}

		var weight *int
		if entity.weight.Valid {
			temp := int(entity.weight.Int64)
			weight = &temp
		}

		list = append(list, domain.Nutrition{
			Date:     date,
			Calories: calories,
			Weight:   weight,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *sqlNutritionRepository) FindAverageNutrition(ctx context.Context, start, end time.Time) (domain.AverageNutrition, error) {
	//TODO implement me
	panic("implement me")
}

func (s *sqlNutritionRepository) Create(ctx context.Context, n domain.Nutrition) (domain.Nutrition, error) {
	//TODO implement me
	panic("implement me")
}

func (s *sqlNutritionRepository) Update(ctx context.Context, n domain.Nutrition) (domain.Nutrition, error) {
	//TODO implement me
	panic("implement me")
}

func (s *sqlNutritionRepository) Delete(ctx context.Context, n domain.Nutrition) error {
	//TODO implement me
	panic("implement me")
}
