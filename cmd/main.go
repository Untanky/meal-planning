package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"meal-planning/database"
	"meal-planning/domain"
	myHttp "meal-planning/http"
	"net/http"
	"os"
	"path"
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
	indexHandler := &indexHandler{templateHandler: tmplHandler, manifest: myManifest}
	mealHandler := &mealHandler{templateHandler: tmplHandler, mealDayRepo: mealDayRepo}

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
	manifest manifest
}

type indexData struct {
	Manifest manifest
	Meals    []domain.MealDay
}

func (h *indexHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.serveTemplate(writer, "index.gohtml", indexData{
		Manifest: h.manifest,
		Meals: []domain.MealDay{
			{time.Date(2024, 6, 29, 0, 0, 0, 0, time.Local), "", "", "Pizza", []string{}},
			{time.Date(2024, 6, 30, 0, 0, 0, 0, time.Local), "", "", "Pasta", []string{}},
			{time.Date(2024, 6, 31, 0, 0, 0, 0, time.Local), "", "", "Burger", []string{}},
		},
	})
}

type mealHandler struct {
	templateHandler
	mealDayRepo database.MealDayRepository
}

func (h *mealHandler) getMeals(writer http.ResponseWriter, request *http.Request) {
	meals, err := h.mealDayRepo.FindByDateRange(context.TODO(), time.Date(2024, 7, 2, 0, 0, 0, 0, time.Local), time.Date(2024, 7, 6, 0, 0, 0, 0, time.Local))
	fmt.Println(meals, err)

	h.serveTemplate(writer, "meal-list", []domain.MealDay{
		{time.Date(2024, 6, 29, 0, 0, 0, 0, time.Local), "", "", "Pizza", []string{}},
		{time.Date(2024, 6, 30, 0, 0, 0, 0, time.Local), "", "", "Pasta", []string{}},
		{time.Date(2024, 6, 31, 0, 0, 0, 0, time.Local), "", "", "Burger", []string{}},
	})
}

func (h *mealHandler) getMealByDate(writer http.ResponseWriter, request *http.Request) {
	h.serveTemplate(writer, "meal-day", domain.MealDay{
		time.Date(2024, 6, 29, 0, 0, 0, 0, time.Local), "", "", "Pizza", []string{},
	})
}

func (h *mealHandler) getMealFormByDate(writer http.ResponseWriter, request *http.Request) {
	h.serveTemplate(writer, "meal-day-form", domain.MealDay{
		time.Date(2024, 6, 29, 0, 0, 0, 0, time.Local), "", "", "Pizza", []string{}},
	)
}

func (h *mealHandler) updateMealByDate(writer http.ResponseWriter, request *http.Request) {
	h.serveTemplate(writer, "meal-day", domain.MealDay{
		time.Date(2024, 6, 29, 0, 0, 0, 0, time.Local), "", "", "Pizza", []string{}},
	)
}
