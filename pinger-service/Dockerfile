FROM golang:1.23-alpine

WORKDIR /pinger-service

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o pinger-service .

CMD ["./pinger-service"]