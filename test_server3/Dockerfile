FROM golang:1.23

WORKDIR /test3

COPY . .

RUN go mod tidy && go build -o test3

EXPOSE 8083

CMD ["./test3"]