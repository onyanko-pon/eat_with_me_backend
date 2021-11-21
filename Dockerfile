FROM golang:1.17 as builder

WORKDIR /
RUN mkdir /app
COPY . /app/

WORKDIR /app

# RUN go mod tidy
RUN go mod download

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
WORKDIR /app/src
RUN go build \
  -o /go/bin/main


FROM alpine:3.12 as runner

COPY --from=builder /go/bin/main /main

CMD /main
