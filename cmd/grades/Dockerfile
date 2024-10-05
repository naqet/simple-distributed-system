# Build
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . . 

RUN go build -o grades-service cmd/grades/main.go

# Run
FROM scratch

COPY --from=build /app/grades-service .

EXPOSE 3000

CMD [ "/grades-service" ]
