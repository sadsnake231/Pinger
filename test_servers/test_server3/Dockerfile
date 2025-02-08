FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o test_server3 .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/test_server3 .

EXPOSE 8083

CMD ["./test_server3"]