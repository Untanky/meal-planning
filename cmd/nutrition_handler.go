package main

import (
	"meal-planning/domain"
	"net/http"
)

type nutritionHandler struct {
	templateHandler
	manifest         manifest
	nutritionService *domain.NutritionService
}

func (h *nutritionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//start := time.Now()
	//end := start.Add(7 * 24 * time.Hour)

	//_, err := h.nutritionService.FindByDateRange(context.TODO(), start, end)
	//if err != nil {
	//	slog.Error("error retrieving meals from repository", slog.Any("reason", err))
	//	http.Error(writer, "failed retrieving meal days", http.StatusInternalServerError)
	//	return
	//}

	h.serveTemplate(writer, "nutrition.gohtml", indexData{
		Manifest: h.manifest,
		Meals:    []domain.MealDay{},
	})
}
