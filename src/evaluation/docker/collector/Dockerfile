FROM golang:1.14
MAINTAINER  Maksim Baldin maksim.baldin@online.liverpool.ac.uk

COPY src/evaluation/cfg/service_a_supervisord.conf /etc/supervisor/conf.d/supervisord.conf

WORKDIR /collector

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    COLLECTOR_CFG=/collector/src/implementation/collector/cfg/collector.yaml


# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls

RUN go build -o collector src/implementation/collector/collector.go

CMD ["./collector"]