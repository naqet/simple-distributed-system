package registry

import (
	"context"
	"distributed-go/utils"
	"log"
	"net/http"
	"sync"
	"time"
)

const PORT = "3001"
const URL = "http://localhost:" + PORT

type registryService struct {
	registry map[string][]string
	mu       sync.RWMutex
}

func New(url string) *registryService {
	return &registryService{map[string][]string{}, sync.RWMutex{}}
}

func (r *registryService) URL() string {
	url, err := utils.GetServiceURL(PORT)

	if err != nil {
		log.Println("Service url could not be obtained", err)
		return "http://localhost:" + PORT
	}

	return url
}

func (r *registryService) Name() string {
	return "Registry"
}

func (r *registryService) Handler() http.Handler {
	go r.checkHealth()
	return r.register()
}

func (reg *registryService) register() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
        name := r.FormValue("name")
        addr := r.FormValue("addr")

		if len(name) == 0 || len(addr) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reg.mu.Lock()
		hosts, ok := reg.registry[string(name)]
        reg.mu.Unlock()
		if !ok {
			reg.registry[string(name)] = []string{addr}
			log.Printf("Registered new service: %s at %s", string(name), addr)
			return
		}

		found := false
		for _, host := range hosts {
			if host == addr {
				found = true
			}
		}

		if found {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("This service is already registered"))
			return
		}

		hosts = append(hosts, addr)
		reg.mu.RLock()
		reg.registry[string(name)] = hosts
		reg.mu.RUnlock()

		log.Printf("Registered new service: %s at %s", string(name), addr)
	})

	mux.HandleFunc("POST /unregister", func(w http.ResponseWriter, r *http.Request) {
        name := r.FormValue("name")
        addr := r.FormValue("addr")

		if len(name) == 0 || len(addr) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reg.mu.Lock()
		hosts, ok := reg.registry[string(name)]
        reg.mu.Unlock()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service is not registered"))
			return
		}

		idx := -1
		for i, host := range hosts {
			if host == addr {
				idx = i
				break
			}
		}

		if idx == -1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Service is not registered"))
			return
		}

        reg.mu.RLock()
		reg.registry[string(name)] = append(hosts[:idx], hosts[idx+1:]...)
        reg.mu.RUnlock()
		log.Printf("Unregistered service: %s at %s", string(name), addr)
	})

	return mux
}

func (r *registryService) checkHealth() {
	ticker := time.NewTicker(3 * time.Second)

	for range ticker.C {
		r.mu.Lock()
		for name, hosts := range r.registry {
			go func() {
				for i, host := range hosts {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()

					req, err := http.NewRequestWithContext(ctx, http.MethodGet, host+ "/health", nil)

					if err != nil {
						log.Printf("Request creation for %s service at %s failed.\n", name, host)
						continue
					}

					res, err := http.DefaultClient.Do(req)

					if err != nil || res.StatusCode != http.StatusOK {
						log.Printf("%s service at %s is not reachable. Updating registry.\n", name, host)
						r.registry[name] = append(hosts[:i], hosts[i+1:]...)
					}
				}
			}()
		}
		r.mu.Unlock()
	}
}
