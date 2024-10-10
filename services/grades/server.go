package grades

import (
	"distributed-go/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const GRADES_SERVICE = "Grades"

type Student struct {
	Name   string  `json:"name"`
	Grades []grade `json:"grades"`
}
func (s Student) GetAverage() int {
    if len(s.Grades) == 0 {
        return 0
    }
	average := 0
	for _, grade := range s.Grades {
		average += grade.Score
	}

	average = average / len(s.Grades)

	return average
}

type grade struct {
	Activity string `json:"activity"`
	Score    int    `json:"score"`
}

type gradesService struct {
	port     string
	students []Student
}

func New(port string) *gradesService {
	return &gradesService{port, []Student{
		{Name: "John Smith", Grades: []grade{}},
		{Name: "Katie Jerry", Grades: []grade{}},
		{Name: "Jacob Holden", Grades: []grade{}},
	}}
}

func (g *gradesService) Name() string {
	return GRADES_SERVICE
}

func (g *gradesService) Port() string {
	return utils.GetPort(g.port)
}

func (g *gradesService) Handler() http.Handler {
	return g.register()
}

func (g *gradesService) register() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /grades", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		activity := r.FormValue("activity")
		score := r.FormValue("score")

		if len(name) == 0 || len(activity) == 0 || len(score) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("None of the form values can be empty."))
			return
		}

		idx := -1
		for i, student := range g.students {
			if student.Name == name {
				idx = i
				break
			}
		}

		if idx == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Student doesn't exist"))
			return
		}

		scoreInt, err := strconv.Atoi(score)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Score is not an integer"))
			return
		}

		g.students[idx].Grades = append(g.students[idx].Grades, grade{activity, scoreInt})
	})

	mux.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
        data, err := json.Marshal(g.students)

        if err != nil {
            log.Println("Problem with marshaling students", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.Header().Add("Content-Type", "application/json")
        w.Write(data)
	})

	mux.HandleFunc("/students/{name}", func(w http.ResponseWriter, r *http.Request) {
        name := r.PathValue("name")
        idx := -1;
        for i, student := range g.students {
            if student.Name == name {
                idx = i;
                break;
            }
        }

        if idx == -1 {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        data, err := json.Marshal(g.students[idx])

        if err != nil {
            log.Println("Problem with marshaling student", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.Header().Add("Content-Type", "application/json")
        w.Write(data)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy"))
	})

	return mux
}
