FROM golang:alpine

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o builded ./cmd

EXPOSE 8080

CMD ["./builded"]
