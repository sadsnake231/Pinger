FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o test_server2 .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/test_server2 .

EXPOSE 8082

CMD ["./test_server2"]