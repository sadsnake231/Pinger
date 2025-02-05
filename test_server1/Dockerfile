FROM golang:1.23

WORKDIR /test1

COPY . .

RUN go mod tidy && go build -o test1

EXPOSE 8081

CMD ["./test1"]