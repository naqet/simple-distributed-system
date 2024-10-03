package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Service interface {
    Name() string
    URL() string
    Handler() http.Handler
}

func Run(ctx context.Context, service Service) context.Context {
    ctx = startService(ctx, service)

    return ctx;
}

func startService(ctx context.Context, service Service) context.Context {
    ctx, cancel := context.WithCancel(ctx);
    var server http.Server
    server.Addr = service.URL()
    server.Handler = service.Handler()

    go func() {
        log.Println(server.ListenAndServe())
        cancel()
    }()

    go func() {
        fmt.Printf("%s service started. Enter any key to stop it\n", service.Name())
        var s string;
        fmt.Scanln(&s)
        server.Shutdown(ctx)
    }()

    return ctx;
}
