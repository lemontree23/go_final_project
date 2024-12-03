FROM golang:1.22.3 AS builder

WORKDIR /scheduler

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o scheduler ./cmd/scheduler/main.go

FROM ubuntu:latest

WORKDIR /scheduler

RUN mkdir storage

COPY --from=builder /scheduler/scheduler /scheduler
COPY ./config config
COPY ./web web

CMD ["./scheduler"]

