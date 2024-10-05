package service

import (
	"context"
	"distributed-go/services/registry"
	"log"
	"net/http"
)

type service interface {
	Name() string
	Port() string
	Handler() http.Handler
}

const SHOULD_REGISTER = "SHOULD_REGISTER"

func Run(ctx context.Context, service service, register bool) context.Context {
	ctx = context.WithValue(ctx, SHOULD_REGISTER, register)
	ctx = startService(ctx, service)

	if register {
		err := registry.RegisterService(service.Name(), service.Port())

		if err != nil {
			log.Printf("%s service could not be registered.\nError: %s\n", service.Name(), err)
		}
	}

	return ctx
}

func startService(ctx context.Context, service service) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	var server http.Server
    server.Addr = ":" + service.Port()
	server.Handler = service.Handler()

	go func() {
		log.Println(server.ListenAndServe())
		val, ok := ctx.Value(SHOULD_REGISTER).(bool)

		if ok && val {
			err := registry.UnregisterService(service.Name(), service.Port())

			if err != nil {
				log.Printf("%s service could not be unregistered.\nError: %s\n", service.Name(), err)
			}
		}

		cancel()
	}()

	return ctx
}
