package logger

import (
	"distributed-go/utils"
	"io"
	"log"
	"net/http"
	"os"
)

type logger string

func (l logger) Write(data []byte) (int, error) {
	file, err := os.OpenFile(string(l), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	return file.Write(data)
}

type loggerService struct {
	port string
}

func New(port string) *loggerService {
	return &loggerService{port}
}

func (l *loggerService) URL() string {
	url, err := utils.GetServiceURL(l.port)

	if err != nil {
		log.Println("Service url could not be obtained", err)
		return "http://localhost:" + l.port
	}

	return url
}

func (l *loggerService) Name() string {
	return "Logger"
}

func (l *loggerService) Handler() http.Handler {
	clog := log.New(logger("./app.log"), "", log.LstdFlags)
	return register(clog)
}

func register(clog *log.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /log", func(w http.ResponseWriter, r *http.Request) {
		msg, err := io.ReadAll(r.Body)

		if err != nil || len(msg) == 0 {
			if err != nil {
				log.Println("Error while reading r.Body: ", err)
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("Received request from:", r.RemoteAddr)

		clog.Println(string(msg))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy"))
	})

	return mux
}
