package registry

import "os"

func getRegistryURL() string {
	url := os.Getenv("REGISTRY_URL")
	if url == "" {
		return "http://localhost:" + DEFAULT_PORT
	}

    return url
}
