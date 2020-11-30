FROM golang:1.14
MAINTAINER  Maksim Baldin maksim.baldin@online.liverpool.ac.uk

RUN apt-get update && apt-get install -y telnet

ARG binary_name
ARG service_config_file

WORKDIR /service

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    CFG=src/evaluation/cfg/$service_config_file

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls

RUN go build -o service_a src/evaluation/cmd/run_service.go

CMD ["./service_a"]