package main

import (
	"context"
	"log/slog"
	"meal-planning/domain"
	"net/http"
	"time"
)

type nutritionHandler struct {
	templateHandler
	manifest         manifest
	nutritionService *domain.NutritionService
}

type nutritionData struct {
	Manifest         manifest
	NutritionEntries []nutritionView
}

type nutritionView struct {
	Date     time.Time
	Calories *int
	Weight   *float64
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
		var weight *float64
		if nutrition.Weight != nil {
			temp := float64(*nutrition.Weight) / 1000
			weight = &temp
		}

		nutritionEntries[i] = nutritionView{
			Date:     nutrition.Date,
			Calories: nutrition.Calories,
			Weight:   weight,
		}
	}

	h.serveTemplate(writer, "nutrition.gohtml", nutritionData{
		Manifest:         h.manifest,
		NutritionEntries: nutritionEntries,
	})
}
