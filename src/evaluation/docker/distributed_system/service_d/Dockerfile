FROM golang:1.14
MAINTAINER  Maksim Baldin maksim.baldin@online.liverpool.ac.uk

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

RUN go build -o service_d src/evaluation/cmd/run_service.go

CMD ["./service_d"]