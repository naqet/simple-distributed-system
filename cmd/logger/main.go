package main

import (
	"context"
	"distributed-go/internal/service"
	"distributed-go/services/logger"
)

func main() {
    logService := logger.New(":3000")
    ctx := service.Run(context.Background(), logService, true)

    <-ctx.Done()
}
