package logger

import (
	"distributed-go/utils"
	"io"
	"log"
	"net/http"
	"os"
)

const LOGGER_SERVICE = "Logger"

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
	defaultPort string
}

func New(defaultPort string) *loggerService {
	return &loggerService{defaultPort}
}

func (l *loggerService) Port() string {
    return utils.GetPort(l.defaultPort)
}

func (l *loggerService) Name() string {
	return LOGGER_SERVICE
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
