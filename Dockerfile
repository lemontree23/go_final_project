FROM golang:1.22 AS builder

WORKDIR /scheduler

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o scheduler ./cmd/scheduler/main.go

FROM ubuntu:latest

WORKDIR /scheduler

COPY --from=builder /scheduler/scheduler /scheduler
COPY ./config config
COPY ./web web
COPY ./storage storage

EXPOSE 7540

CMD ["./scheduler"]

