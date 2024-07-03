package main

import (
	"encoding/json"
	"html/template"
	"log/slog"
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

type mealDay struct {
	Date      time.Time
	Breakfast string
	Lunch     string
	Dinner    string
	Snacks    []string
}

var globalManifest manifest

func main() {
	slog.Info("Starting application")

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

	globalManifest = manifest{}
	for _, file := range mFile.OutputFiles {
		if path.Ext(file) == ".css" {
			globalManifest.CssFiles = append(globalManifest.CssFiles, file)
		} else if path.Ext(file) == ".js" {
			globalManifest.JsFiles = append(globalManifest.JsFiles, file)
		} else {
			slog.Warn("Unknown file extension", slog.Any("file", file))
		}
	}

	tmpl, err := template.ParseGlob("./views/*.gohtml")
	if err != nil {
		slog.Error("Error parsing template", slog.Any("reason", err))
		panic(err)
	}

	indexHandler := &handleIndex{template: tmpl}
	mealHandler := &mealHandler{template: tmpl}

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

type handleIndex struct {
	template *template.Template
}

type indexData struct {
	Manifest manifest
	Meals    []mealDay
}

func (h *handleIndex) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	myWriter := myHttp.NewBufferedResponseWriter(writer)

	err := h.template.ExecuteTemplate(myWriter, "index.gohtml", indexData{
		Manifest: globalManifest,
		Meals: []mealDay{
			{time.Date(2024, 6, 29, 0, 0, 0, 0, time.Local), "", "", "Pizza", []string{}},
			{time.Date(2024, 6, 30, 0, 0, 0, 0, time.Local), "", "", "Pasta", []string{}},
			{time.Date(2024, 6, 31, 0, 0, 0, 0, time.Local), "", "", "Burger", []string{}},
		},
	})
	if err != nil {
		http.Error(writer, "could not render template", http.StatusInternalServerError)
	}

	header := writer.Header()
	header.Add("Content-Type", "text/html; charset=utf-8")
	header.Add("Cache-Control", "no-store")

	err = myWriter.Close()
	if err != nil {
		slog.Error("error writing to response", slog.Any("reason", err))
		panic(err)
	}
}

type mealHandler struct {
	template *template.Template
}

func (h *mealHandler) getMeals(writer http.ResponseWriter, request *http.Request) {
	myWriter := myHttp.NewBufferedResponseWriter(writer)

	err := h.template.ExecuteTemplate(myWriter, "meal-list", []mealDay{
		{time.Date(2024, 6, 29, 0, 0, 0, 0, time.Local), "", "", "Pizza", []string{}},
		{time.Date(2024, 6, 30, 0, 0, 0, 0, time.Local), "", "", "Pasta", []string{}},
		{time.Date(2024, 6, 31, 0, 0, 0, 0, time.Local), "", "", "Burger", []string{}},
	})
	if err != nil {
		http.Error(writer, "could not render template", http.StatusInternalServerError)
	}

	header := writer.Header()
	header.Add("Content-Type", "text/html; charset=utf-8")
	header.Add("Cache-Control", "no-store")

	err = myWriter.Close()
	if err != nil {
		slog.Error("error writing to response", slog.Any("reason", err))
		panic(err)
	}
}

func (h *mealHandler) getMealByDate(writer http.ResponseWriter, request *http.Request) {
	myWriter := myHttp.NewBufferedResponseWriter(writer)

	err := h.template.ExecuteTemplate(myWriter, "meal-day", mealDay{
		time.Date(2024, 6, 29, 0, 0, 0, 0, time.Local), "", "", "Pizza", []string{}},
	)
	if err != nil {
		http.Error(writer, "could not render template", http.StatusInternalServerError)
	}

	header := writer.Header()
	header.Add("Content-Type", "text/html; charset=utf-8")
	header.Add("Cache-Control", "no-store")

	err = myWriter.Close()
	if err != nil {
		slog.Error("error writing to response", slog.Any("reason", err))
		panic(err)
	}
}

func (h *mealHandler) getMealFormByDate(writer http.ResponseWriter, request *http.Request) {
	myWriter := myHttp.NewBufferedResponseWriter(writer)

	err := h.template.ExecuteTemplate(myWriter, "meal-day-form", mealDay{
		time.Date(2024, 6, 29, 0, 0, 0, 0, time.Local), "", "", "Pizza", []string{}},
	)
	if err != nil {
		http.Error(writer, "could not render template", http.StatusInternalServerError)
	}

	header := writer.Header()
	header.Add("Content-Type", "text/html; charset=utf-8")
	header.Add("Cache-Control", "no-store")

	err = myWriter.Close()
	if err != nil {
		slog.Error("error writing to response", slog.Any("reason", err))
		panic(err)
	}
}

func (h *mealHandler) updateMealByDate(writer http.ResponseWriter, request *http.Request) {
	myWriter := myHttp.NewBufferedResponseWriter(writer)

	err := h.template.ExecuteTemplate(myWriter, "meal-day", mealDay{
		time.Date(2024, 6, 29, 0, 0, 0, 0, time.Local), "", "", "Pizza", []string{}},
	)
	if err != nil {
		http.Error(writer, "could not render template", http.StatusInternalServerError)
	}

	header := writer.Header()
	header.Add("Content-Type", "text/html; charset=utf-8")
	header.Add("Cache-Control", "no-store")

	err = myWriter.Close()
	if err != nil {
		slog.Error("error writing to response", slog.Any("reason", err))
		panic(err)
	}
}
