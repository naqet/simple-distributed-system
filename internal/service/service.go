package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func Run(ctx context.Context, name, url string, handler http.Handler) context.Context {
    ctx = startService(ctx, name, url, handler)

    return ctx;
}

func startService(ctx context.Context, name, url string, handler http.Handler) context.Context {
    ctx, cancel := context.WithCancel(ctx);
    var server http.Server
    server.Addr = url
    server.Handler = handler

    go func() {
        log.Println(server.ListenAndServe())
        cancel()
    }()

    go func() {
        fmt.Printf("%s service started. Enter any key to stop it\n", name)
        var s string;
        fmt.Scanln(&s)
        server.Shutdown(ctx)
    }()

    return ctx;
}
