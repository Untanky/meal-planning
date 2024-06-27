package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

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

func (h *handleIndex) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	buffer := new(bytes.Buffer)
	hash := md5.New()
	multi := io.MultiWriter(hash, buffer)

	err := h.template.ExecuteTemplate(multi, "index.gohtml", nil)
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
