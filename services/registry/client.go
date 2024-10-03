package registry

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func RegisterService(name string) error {
	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Post(URL, "text/plain", strings.NewReader(name))

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Request failed with status %d.\n", res.StatusCode)
	}

	return nil
}

func UnregisterService(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, URL, strings.NewReader(name))

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Request failed with status %d.\n", res.StatusCode)
	}

	return nil
}
