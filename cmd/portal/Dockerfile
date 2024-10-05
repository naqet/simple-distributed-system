# Build
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . . 

RUN go build -o portal-service cmd/portal/main.go

# Run
FROM scratch

COPY --from=build /app/portal-service .
COPY --from=build /app/services/portal/templates/*.html /services/portal/templates/

EXPOSE 3000

CMD [ "/portal-service" ]
