package main

import (
	"context"
	"encoding/json"
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
	mealDayService := domain.NewMealDayService(mealDayRepo)

	nutritionRepo := database.NewSqlNutritionRepository(db)
	nutritionService := domain.NewNutritionService(nutritionRepo)

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
		mealDayService:  mealDayService,
	}

	nutritionHandler := &nutritionHandler{
		templateHandler:  tmplHandler,
		manifest:         myManifest,
		nutritionService: nutritionService,
	}

	mealHandler := &mealHandler{
		templateHandler: tmplHandler,
		mealDayService:  mealDayService,
	}

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("GET /meals", mealHandler.getMeals)
	mux.HandleFunc("GET /meals/{date}", mealHandler.getMealByDate)
	mux.HandleFunc("PUT /meals/{date}", mealHandler.updateMealByDate)
	mux.HandleFunc("GET /meals/{date}/form", mealHandler.getMealFormByDate)
	mux.Handle("/", indexHandler)
	mux.Handle("/nutrition", nutritionHandler)

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
		slog.Error("Error executing template", slog.Any("reason", err))
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
	manifest       manifest
	mealDayService *domain.MealDayService
}

type indexData struct {
	Manifest manifest
	Meals    []domain.MealDay
}

func (h *indexHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	end := start.Add(7 * 24 * time.Hour)

	meals, err := h.mealDayService.FindByDateRange(context.TODO(), start, end)
	if err != nil {
		slog.Error("error retrieving meals from repository", slog.Any("reason", err))
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
	mealDayService *domain.MealDayService
}

func (h *mealHandler) getMeals(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	end := start.Add(7 * 24 * time.Hour)

	meals, err := h.mealDayService.FindByDateRange(context.TODO(), start, end)
	if err != nil {
		slog.Error("error retrieving meals from repository", slog.Any("reason", err))
		http.Error(writer, "failed retrieving meal days", http.StatusInternalServerError)
		return
	}

	h.serveTemplate(writer, "meal-list", meals)
}

func (h *mealHandler) getMealByDate(writer http.ResponseWriter, request *http.Request) {
	dateString := request.PathValue("date")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		slog.Error("error parsing date", slog.Any("reason", err))
		http.Error(writer, "date must be an ISO date", http.StatusBadRequest)
		return
	}

	meal, err := h.mealDayService.FindByDate(context.TODO(), date)
	if err != nil {
		slog.Error("error retrieving meal from service", slog.Any("reason", err))
		http.Error(writer, "failed retrieving meal", http.StatusInternalServerError)
		return
	}

	h.serveTemplate(writer, "meal-day", meal)
}

func (h *mealHandler) getMealFormByDate(writer http.ResponseWriter, request *http.Request) {
	dateString := request.PathValue("date")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		slog.Error("error parsing date", slog.Any("reason", err))
		http.Error(writer, "date must be an ISO date", http.StatusBadRequest)
		return
	}

	meal, err := h.mealDayService.FindByDate(context.TODO(), date)
	if err != nil {
		slog.Error("error retrieving meal from service", slog.Any("reason", err))
		http.Error(writer, "failed retrieving meal", http.StatusInternalServerError)
		return
	}

	h.serveTemplate(writer, "meal-day-form", meal)
}

func (h *mealHandler) updateMealByDate(writer http.ResponseWriter, request *http.Request) {
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

	meal := domain.MealDay{
		Date:      date,
		Breakfast: request.Form.Get("breakfast"),
		Lunch:     request.Form.Get("lunch"),
		Dinner:    request.Form.Get("dinner"),
		Snacks:    strings.Split(request.Form.Get("snacks"), ","),
	}

	meal, err = h.mealDayService.Upsert(context.TODO(), meal)
	if err != nil {
		slog.Error("error updating meal", slog.Any("reason", err))
		http.Error(writer, "failed updating meal", http.StatusInternalServerError)
		return
	}

	h.serveTemplate(writer, "meal-day", meal)
}
