package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"meal-planning/domain"
	"net/http"
	"strconv"
	"time"
)

type nutritionHandler struct {
	templateHandler
	manifest         manifest
	nutritionService *domain.NutritionService
}

type nutritionData struct {
	Manifest                    manifest
	NutritionEntries            []nutritionView
	NutritionJSON               string
	TotalDailyEnergyExpenditure totalDailyEnergyExpenditureView
}

type totalDailyEnergyExpenditureView struct {
	Start                       time.Time
	End                         time.Time
	AverageCalories             int
	PeriodWeightDifference      float64
	TotalDailyEnergyExpenditure int
}

type nutritionView struct {
	Date     time.Time `json:"date"`
	Calories int       `json:"calories,omitempty"`
	Weight   float64   `json:"weight,omitempty"`
}

func (h *nutritionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)

	nutritionList, err := h.nutritionService.FindByDateRange(context.TODO(), start, end)
	if err != nil {
		slog.Error("error retrieving meals from repository", slog.Any("reason", err))
		http.Error(writer, "failed retrieving meal days", http.StatusInternalServerError)
		return
	}

	nutritionEntries := make([]nutritionView, len(nutritionList))
	for i, nutrition := range nutritionList {

		nutritionEntries[i] = nutritionView{
			Date:     nutrition.Date,
			Calories: nutrition.Calories,
			Weight:   float64(nutrition.Weight) / 1000,
		}
	}

	nutritionJSON, err := json.Marshal(nutritionEntries)
	if err != nil {
		slog.Error("error marshaling meals to JSON", slog.Any("reason", err))
		http.Error(writer, "failed marshaling meals to JSON", http.StatusInternalServerError)
		return
	}

	totalDailyEnergyExpenditure, err := h.nutritionService.CalculateTotalDailyEnergyExpenditure(context.TODO(), start, end)
	if err != nil {
		slog.Error("error retrieving meals from repository", slog.Any("reason", err))
		http.Error(writer, "failed retrieving meal days", http.StatusInternalServerError)
		return
	}

	h.serveTemplate(writer, "nutrition.gohtml", nutritionData{
		Manifest:         h.manifest,
		NutritionEntries: nutritionEntries,
		NutritionJSON:    string(nutritionJSON),
		TotalDailyEnergyExpenditure: totalDailyEnergyExpenditureView{
			Start:                       totalDailyEnergyExpenditure.Start,
			End:                         totalDailyEnergyExpenditure.End,
			AverageCalories:             totalDailyEnergyExpenditure.AverageCalories,
			PeriodWeightDifference:      float64(totalDailyEnergyExpenditure.PeriodWeightDifference) / 1000,
			TotalDailyEnergyExpenditure: totalDailyEnergyExpenditure.TotalDailyEnergyExpenditure,
		},
	})
}

func (h *nutritionHandler) updateNutritionEntry(writer http.ResponseWriter, request *http.Request) {
	dateString := request.PathValue("date")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		slog.Error("error parsing date", slog.Any("reason", err))
		http.Error(writer, "date must be an ISO date", http.StatusBadRequest)
		return
	}

	err = request.ParseForm()
	if err != nil {
		slog.Error("error parsing form", slog.Any("reason", err))
		http.Error(writer, "could not parse form", http.StatusBadRequest)
		return
	}

	calories, err := strconv.Atoi(request.FormValue("calories"))
	if err == nil {
		slog.Warn("error parsing calories", slog.Any("reason", err))
	}

	weight, err := strconv.ParseFloat(request.FormValue("weight"), 64)
	if err != nil {
		slog.Warn("error parsing weight", slog.Any("reason", err))
	}

	nutrition := domain.Nutrition{
		Date:     date,
		Calories: calories,
		Weight:   int(weight * 1000),
	}

	nutrition, err = h.nutritionService.Upsert(context.TODO(), nutrition)
	if err != nil {
		slog.Error("error updating nutrition", slog.Any("reason", err))
		http.Error(writer, "could not update nutrition", http.StatusInternalServerError)
		return
	}

	nutritionEntry := nutritionView{
		Date:     nutrition.Date,
		Calories: nutrition.Calories,
		Weight:   float64(nutrition.Weight) / 1000,
	}

	nutritionJSON, err := json.Marshal(nutritionEntry)
	if err != nil {
		slog.Error("error marshaling meals to JSON", slog.Any("reason", err))
		http.Error(writer, "failed marshaling meals to JSON", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("HX-Trigger", fmt.Sprintf(`{ "updateNutritionData": %s }`, string(nutritionJSON)))

	h.serveTemplate(writer, "nutrition-entry", nutritionEntry)
}
