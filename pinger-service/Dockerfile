FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o pinger-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pinger-service .

CMD ["./pinger-service"]