package main

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"log/slog"
	"meal-planning/database"
	"meal-planning/domain"
	myHttp "meal-planning/http"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type manifestFile struct {
	OutputFiles []string
}

type manifest struct {
	CssFiles []string
	JsFiles  []string
}

func main() {
	slog.Info("Starting application")

	slog.Info("Connecting to database")
	db := connectDatabase()
	defer db.Close()
	migrateDatabase(db)

	mealDayRepo := database.NewSqlMealDayRepository(db)

	slog.Info("Loading manifest")
	file, err := os.OpenFile("./manifest.json", os.O_RDONLY, os.ModePerm)
	if err != nil {
		slog.Error("Failed to open manifest.json", slog.Any("reason", err))
		panic(err)
	}

	mFile := manifestFile{}
	err = json.NewDecoder(file).Decode(&mFile)
	if err != nil {
		slog.Error("Failed to parse manifest.json", slog.Any("reason", err))
		panic(err)
	}

	myManifest := manifest{}
	for _, file := range mFile.OutputFiles {
		if path.Ext(file) == ".css" {
			myManifest.CssFiles = append(myManifest.CssFiles, file)
		} else if path.Ext(file) == ".js" {
			myManifest.JsFiles = append(myManifest.JsFiles, file)
		} else {
			slog.Warn("Unknown file extension", slog.Any("file", file))
		}
	}

	tmpl, err := template.ParseGlob("./views/*.gohtml")
	if err != nil {
		slog.Error("Error parsing template", slog.Any("reason", err))
		panic(err)
	}

	tmplHandler := templateHandler{
		template: tmpl,
	}
	indexHandler := &indexHandler{
		templateHandler: tmplHandler,
		manifest:        myManifest,
		mealDayRepo:     mealDayRepo,
	}
	mealHandler := &mealHandler{
		templateHandler: tmplHandler,
		mealDayRepo:     mealDayRepo,
	}

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("GET /meals", mealHandler.getMeals)
	mux.HandleFunc("GET /meals/{date}", mealHandler.getMealByDate)
	mux.HandleFunc("PUT /meals/{date}", mealHandler.updateMealByDate)
	mux.HandleFunc("GET /meals/{date}/form", mealHandler.getMealFormByDate)
	mux.Handle("/", indexHandler)

	slog.Info("Starting server")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("error running server", slog.Any("reason", err))
		panic(err)
	}
}

type templateHandler struct {
	template *template.Template
}

func (handler *templateHandler) serveTemplate(writer http.ResponseWriter, name string, data interface{}) {
	bufferedWriter := myHttp.NewBufferedResponseWriter(writer)

	err := handler.template.ExecuteTemplate(bufferedWriter, name, data)
	if err != nil {
		http.Error(writer, "could not render template", http.StatusInternalServerError)
		return
	}

	header := writer.Header()
	header.Add("Content-Type", "text/html; charset=utf-8")
	header.Add("Cache-Control", "no-store")

	err = bufferedWriter.Close()
	if err != nil {
		slog.Error("error writing response", slog.Any("reason", err))
		panic(err)
	}
}

type indexHandler struct {
	templateHandler
	manifest    manifest
	mealDayRepo database.MealDayRepository
}

type indexData struct {
	Manifest manifest
	Meals    []domain.MealDay
}

func (h *indexHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	meals, err := h.mealDayRepo.FindByDateRange(context.TODO(), time.Date(2024, 7, 2, 0, 0, 0, 0, time.Local), time.Date(2024, 7, 6, 0, 0, 0, 0, time.Local))
	if err != nil {
		http.Error(writer, "failed retrieving meal days", http.StatusInternalServerError)
		return
	}

	h.serveTemplate(writer, "index.gohtml", indexData{
		Manifest: h.manifest,
		Meals:    meals,
	})
}

type mealHandler struct {
	templateHandler
	mealDayRepo database.MealDayRepository
}

func (h *mealHandler) getMeals(writer http.ResponseWriter, request *http.Request) {
	meals, err := h.mealDayRepo.FindByDateRange(context.TODO(), time.Date(2024, 7, 2, 0, 0, 0, 0, time.Local), time.Date(2024, 7, 6, 0, 0, 0, 0, time.Local))
	if err != nil {
		http.Error(writer, "failed retrieving meal days", http.StatusInternalServerError)
		return
	}

	h.serveTemplate(writer, "meal-list", meals)
}

func (h *mealHandler) getMealByDate(writer http.ResponseWriter, request *http.Request) {
	dateString := request.PathValue("date")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		http.Error(writer, "date must be an ISO date", http.StatusBadRequest)
		return
	}

	meal, err := h.mealDayRepo.FindByDate(context.TODO(), date)
	if err != nil {
		http.Error(writer, "failed retrieving meal days", http.StatusInternalServerError)
		return
	}

	h.serveTemplate(writer, "meal-day", meal)
}

func (h *mealHandler) getMealFormByDate(writer http.ResponseWriter, request *http.Request) {
	dateString := request.PathValue("date")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		http.Error(writer, "date must be an ISO date", http.StatusBadRequest)
		return
	}

	meal, err := h.mealDayRepo.FindByDate(context.TODO(), date)
	if err != nil {
		http.Error(writer, "failed retrieving meal days", http.StatusInternalServerError)
		return
	}

	h.serveTemplate(writer, "meal-day-form", meal)
}

func (h *mealHandler) updateMealByDate(writer http.ResponseWriter, request *http.Request) {
	dateString := request.PathValue("date")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		http.Error(writer, "date must be an ISO date", http.StatusBadRequest)
		return
	}

	err = request.ParseForm()
	if err != nil {
		http.Error(writer, "could not parse form", http.StatusBadRequest)
		return
	}

	meal, err := h.mealDayRepo.FindByDate(context.TODO(), date)
	if err != nil && !errors.Is(err, database.NotFound) {
		slog.Error("error retrieving meal days", slog.Any("reason", err))
		http.Error(writer, "failed retrieving meal days", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, database.NotFound) {
		meal, err = h.mealDayRepo.Create(context.TODO(), domain.MealDay{
			Date:      date,
			Breakfast: request.Form.Get("breakfast"),
			Lunch:     request.Form.Get("lunch"),
			Dinner:    request.Form.Get("dinner"),
			Snacks:    strings.Split(request.Form.Get("snacks"), ","),
		})
	} else {
		meal.Breakfast = request.Form.Get("breakfast")
		meal.Lunch = request.Form.Get("lunch")
		meal.Dinner = request.Form.Get("dinner")
		meal.Snacks = strings.Split(request.Form.Get("snacks"), ",")
		meal, err = h.mealDayRepo.Update(context.TODO(), meal)
	}

	h.serveTemplate(writer, "meal-day", meal)
}
