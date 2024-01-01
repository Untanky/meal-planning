package meals

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

type MealPlanningHandler struct {
	indexFile *os.File
}

func NewMealPlanningHandler() (http.Handler, error) {
	file, err := os.Open("index.html")
	if err != nil {
		return nil, err
	}

	handler := &MealPlanningHandler{
		indexFile: file,
	}

	return handler, err
}

func (m *MealPlanningHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.URL.Path)
	switch {
	case (request.URL.Path == "/" || request.URL.Path == "/index.html") && request.Method == http.MethodGet:
		m.serveIndexFile(writer, request)
	case request.URL.Path == "/days" && request.Method == http.MethodGet:
		m.serveDays(writer, request)
	}
}

func (m *MealPlanningHandler) serveIndexFile(writer http.ResponseWriter, request *http.Request) {
	_, err := io.Copy(writer, m.indexFile)
	if err != nil {
		panic(err)
	}

	_, err = m.indexFile.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
}

type day struct {
	DayOfTheWeek string
	Date         string
}

type daysData struct {
	Days []day
}

func (m *MealPlanningHandler) serveDays(writer http.ResponseWriter, request *http.Request) {
	temp, err := template.ParseFiles("days.template.html")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	temp.Execute(
		writer, &daysData{
			Days: []day{
				{"Monday", "01.01.2024"},
				{"Tuesday", "02.01.2024"},
				{"Wednesday", "03.01.2024"},
				{"Thursday", "04.01.2024"},
				{"Friday", "05.01.2024"},
				{"Saturday", "06.01.2024"},
				{"Sunday", "07.01.2024"},
			},
		},
	)

	if err != nil {
		panic(err)
	}
}
