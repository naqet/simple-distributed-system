package main

import (
	"context"
	"distributed-go/internal/service"
	"distributed-go/services/registry"
)

func main() {
    registryService := registry.New()
    ctx := service.Run(context.Background(), registryService, false)

    <-ctx.Done()
}
