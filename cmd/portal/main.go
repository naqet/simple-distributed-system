package main

import (
	"context"
	"distributed-go/internal/service"
	"distributed-go/services/portal"
)

func main() {
	portalService := portal.New("8000")
    ctx := service.Run(context.Background(), portalService, true)

    <- ctx.Done();
}
