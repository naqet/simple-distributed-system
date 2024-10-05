package grades

import (
	"distributed-go/utils"
	"net/http"
	"strconv"
)

type student struct {
	name   string
	grades []grade
}

func (s *student) getAverage() int {
	average := 0
	for _, grade := range s.grades {
		average += grade.score
	}

	average = average / len(s.grades)

	return average
}

type grade struct {
	activity string
	score    int
}

type gradesService struct {
	port     string
	students []student
}

func New(port string) *gradesService {
	return &gradesService{port, []student{
		{name: "John Smith", grades: []grade{}},
		{name: "Katie Jerry", grades: []grade{}},
		{name: "Jacob Holden", grades: []grade{}},
	}}
}

func (g *gradesService) Name() string {
	return "Grades"
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
			if student.name == name {
				idx = i
				break
			}
		}

		if idx == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Student doesn't exist"))
			return
		}

		found := false
		for _, grade := range g.students[idx].grades {
			if grade.activity == activity {
				found = true
				break
			}
		}

		if found {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Activity with such name already exists in the student's grade book."))
			return
		}

		scoreInt, err := strconv.Atoi(score)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Score is not an integer"))
			return
		}

		g.students[idx].grades = append(g.students[idx].grades, grade{activity, scoreInt})
	})

    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Healthy"))
    })

	return mux
}
