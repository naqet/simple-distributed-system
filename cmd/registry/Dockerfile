# Build
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . . 

RUN go build -o registry-service cmd/registry/main.go

# Run
FROM scratch

COPY --from=build /app/registry-service .

EXPOSE 3000

CMD [ "/registry-service" ]
