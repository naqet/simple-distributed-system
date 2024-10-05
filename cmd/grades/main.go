package main;

import (
	"context"
	"distributed-go/internal/service"
	"distributed-go/services/grades"
)

func main() {
    gradesService := grades.New("3002")
    ctx := service.Run(context.Background(), gradesService, true)

    <-ctx.Done()
}
