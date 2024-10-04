package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RegisterService(name, serviceUrl string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    formData := url.Values{}
    formData.Set("name", name)
    formData.Set("addr", serviceUrl)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL + "/register", strings.NewReader(formData.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Request failed with status %d.\n", res.StatusCode)
	}

	return nil
}

func UnregisterService(name, serviceUrl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

    formData := url.Values{}
    formData.Set("name", name)
    formData.Set("addr", serviceUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL + "/unregister", strings.NewReader(formData.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
