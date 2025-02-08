FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o test_server1 .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/test_server1 .

EXPOSE 8081

CMD ["./test_server1"]