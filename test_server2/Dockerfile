FROM golang:1.23

WORKDIR /test2

COPY . .

RUN go mod tidy && go build -o test2

EXPOSE 8082

CMD ["./test2"]