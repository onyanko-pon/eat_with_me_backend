FROM golang:latest

WORKDIR /
RUN mkdir /app
COPY . /app/

WORKDIR /app

RUN go mod tidy

CMD go run /app/src/main.go