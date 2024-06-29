package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type mealDay struct {
	Date      time.Time
	Breakfast string
	Lunch     string
	Dinner    string
	Snacks    []string
}

func main() {
	slog.Info("Starting application")

	tmpl, err := template.ParseGlob("./views/*.gohtml")
	if err != nil {
		slog.Error("Error parsing template", slog.Any("reason", err))
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", &handleIndex{template: tmpl})

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
	Meals []mealDay
}

func (h *handleIndex) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	buffer := new(bytes.Buffer)
	hash := md5.New()
	multi := io.MultiWriter(hash, buffer)

	err := h.template.ExecuteTemplate(multi, "index.gohtml", indexData{
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
	header.Add("Content-Length", fmt.Sprint(buffer.Len()))
	header.Add("Cache-Control", "no-store")

	_, err = io.Copy(writer, buffer)
	if err != nil {
		slog.Error("error writing to response", slog.Any("reason", err))
		panic(err)
	}
}
