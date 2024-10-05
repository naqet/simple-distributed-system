# Build
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . . 

RUN go build -o logger-service cmd/logger/main.go

# Run
FROM scratch

COPY --from=build /app/logger-service .

EXPOSE 3000

CMD [ "/logger-service" ]
