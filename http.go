package meals

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

type MealPlanningHandler struct {
	indexFile *os.File
	templates *template.Template
	db        *leveldb.DB
}

func NewMealPlanningHandler(db *leveldb.DB) (http.Handler, error) {
	file, err := os.Open("index.html")
	if err != nil {
		return nil, err
	}

	templates, err := template.ParseFiles(
		"days.template.html",
	)
	if err != nil {
		return nil, err
	}

	handler := &MealPlanningHandler{
		indexFile: file,
		templates: templates,
		db:        db,
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
	case request.URL.Path == "/days" && request.Method == http.MethodPut:
		m.putDay(writer, request)
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
	Dish         string
}

type daysData struct {
	Days []day
}

func (m *MealPlanningHandler) serveDays(writer http.ResponseWriter, request *http.Request) {
	t := time.Now()
	dayDuration := 24 * time.Hour
	days := make([]day, 7)
	for i, _ := range days {
		dateQuery := t.Format("2006-01-02")
		result, err := m.db.Get([]byte(dateQuery), nil)
		if err == nil {
			days[i] = day{
				t.Format("Monday"),
				t.Format("02.01.2006"),
				string(result),
			}
		} else {
			days[i] = day{
				t.Format("Monday"),
				t.Format("02.01.2006"),
				"",
			}
		}

		t = t.Add(dayDuration)
	}

	err := m.templates.Lookup("days.template.html").Execute(
		writer, &daysData{
			Days: days,
		},
	)
	if err != nil {
		panic(err)
	}
}

func (m *MealPlanningHandler) putDay(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}

	_, err = time.Parse("2006-01-02", request.Form.Get("date"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	matched, err := regexp.Match("[a-zA-Z0-9\\s\\-]*", []byte(request.Form.Get("dinner")))
	if err != nil || !matched {
		writer.WriteHeader(http.StatusBadRequest)
	}

	err = m.db.Put([]byte(request.Form.Get("date")), []byte(request.Form.Get("dinner")), nil)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	m.serveDays(writer, request)
}
