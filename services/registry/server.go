package registry

import (
	"io"
	"log"
	"net/http"
)

const PORT = ":3001"
const URL = "http://localhost" + PORT + "/register"

type registryService struct {
	registry map[string][]string
}

func New(url string) *registryService {
	return &registryService{map[string][]string{}}
}

func (r *registryService) URL() string {
	return PORT
}

func (r *registryService) Name() string {
	return "Registry"
}

func (r *registryService) Handler() http.Handler {
	return r.register()
}

func (reg *registryService) register() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		name, err := io.ReadAll(r.Body)

		if err != nil || len(name) == 0 {
			if err != nil {
				log.Println("Error while reading body", err)
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hosts, ok := reg.registry[string(name)]
		if !ok {
			reg.registry[string(name)] = []string{r.RemoteAddr}
			log.Printf("Registered new service: %s at %s", string(name), r.RemoteAddr)
			return
		}

		found := false
		for _, host := range hosts {
			if host == r.RemoteAddr {
				found = true
			}
		}

		if found {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("This service is already registered"))
			return
		}

		hosts = append(hosts, r.RemoteAddr)
		reg.registry[string(name)] = hosts
		log.Printf("Registered new service: %s at %s", string(name), r.RemoteAddr)
	})

	mux.HandleFunc("DELETE /register", func(w http.ResponseWriter, r *http.Request) {
		name, err := io.ReadAll(r.Body)

		if err != nil || len(name) == 0 {
			if err != nil {
				log.Println("Error while reading body", err)
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hosts, ok := reg.registry[string(name)]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service is not registered"))
			return
		}

		idx := -1
		for i, host := range hosts {
			if host == r.RemoteAddr {
				idx = i
				break
			}
		}

		if idx == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service is not registered"))
			return
		}

		reg.registry[string(name)] = append(hosts[:idx], hosts[idx+1:]...)
		log.Printf("Unregistered service: %s at %s", string(name), r.RemoteAddr)
	})

	return mux
}
