package portal

import (
	"context"
	"distributed-go/services/grades"
	"distributed-go/services/registry"
	"distributed-go/utils"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const PORTAL_SERVICE = "Portal"

type portalService struct {
	port string
}

func New(defaultPort string) *portalService {
	return &portalService{defaultPort}
}

func (p *portalService) Port() string {
	return utils.GetPort(p.port)
}

func (p *portalService) Name() string {
	return PORTAL_SERVICE
}

func (p *portalService) Handler() http.Handler {
	mux := http.NewServeMux()

	templates := template.Must(template.ParseGlob("services/portal/templates/*.html"))

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
        serviceURL, err := registry.GetProvider(grades.GRADES_SERVICE)

        if err != nil {
            log.Println("Error connecting to grades service", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
        defer cancel()

        req, err := http.NewRequestWithContext(ctx, http.MethodGet, serviceURL + "/students", nil)
        
        if err != nil {
            log.Println("Error while creating request", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        res, err := http.DefaultClient.Do(req)

        if err != nil || res.StatusCode != http.StatusOK {
            log.Println("Error getting students", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        
        body, err := io.ReadAll(res.Body)

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        data := []grades.Student{}

        err = json.Unmarshal(body, &data)

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        templates.ExecuteTemplate(w, "index.html", data)
	})
    
    mux.HandleFunc("/student/{name}", func(w http.ResponseWriter, r *http.Request) {
        name := r.PathValue("name")
        
        serviceURL, err := registry.GetProvider(grades.GRADES_SERVICE)

        if err != nil {
            log.Println("Error connecting to grades service", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        parsedURL, err := url.ParseRequestURI(serviceURL)

        if err != nil {
            log.Println("Error while parsing service url", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        finalURL := parsedURL.JoinPath("/students", "/" + name)

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second);
        defer cancel()

        req, err := http.NewRequestWithContext(ctx, http.MethodGet, finalURL.String(), nil)

        if err != nil {
            log.Println("Error while creating request", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        res, err := http.DefaultClient.Do(req)

        if err != nil || res.StatusCode != http.StatusOK {
            log.Println("Error getting student info", err, res.StatusCode)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        
        body, err := io.ReadAll(res.Body)

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        data := grades.Student{}

        err = json.Unmarshal(body, &data)

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        templates.ExecuteTemplate(w, "student.html", data)
    })

    mux.HandleFunc("POST /grade", func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()

        if err != nil {
            log.Println("Error parsing form", err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        serviceURL, err := registry.GetProvider(grades.GRADES_SERVICE)

        if err != nil {
            log.Println("Error connecting to grades service", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second);
        defer cancel()

        req, err := http.NewRequestWithContext(ctx, http.MethodPost, serviceURL + "/grades", strings.NewReader(r.Form.Encode()))

        if err != nil {
            log.Println("Error while creating request", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

        res, err := http.DefaultClient.Do(req)

        if err != nil || res.StatusCode != http.StatusOK {
            log.Println("Error submitting grade info", err, res.StatusCode)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        w.Header().Add("Location", "/student/" + r.Form.Get("name"))
        w.WriteHeader(http.StatusPermanentRedirect)
    })

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy."))
	})

	return mux
}
