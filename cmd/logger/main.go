package main

import (
	"context"
	"distributed-go/internal/service"
	"distributed-go/services/logger"
)

func main() {
    logService := logger.Prepare()
    ctx := service.Run(context.Background(), "Logger", ":3000", logService)

    <-ctx.Done()
}
