package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RegisterService(name, port string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    formData := url.Values{}
    formData.Set("name", name)
    formData.Set("port", port)

    registryUrl := getRegistryURL()

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, registryUrl + "/register", strings.NewReader(formData.Encode()))
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

func UnregisterService(name, port string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

    formData := url.Values{}
    formData.Set("name", name)
    formData.Set("port", port)

    registryUrl := getRegistryURL()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, registryUrl + "/unregister", strings.NewReader(formData.Encode()))
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


func GetProvider(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

    registryURL := getRegistryURL()
    parsedURL, err := url.ParseRequestURI(registryURL)

    if err != nil {
        return "", err
    }

    query := url.Values{}
    query.Add("name", name)
    parsedURL.RawQuery = query.Encode()
    finalURL := parsedURL.JoinPath("/get-provider")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, finalURL.String(), nil)

	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Request failed with status %d.\n", res.StatusCode)
	}
    body, err := io.ReadAll(res.Body)

    if err != nil {
        return "", err
    }

	return string(body), nil
}
