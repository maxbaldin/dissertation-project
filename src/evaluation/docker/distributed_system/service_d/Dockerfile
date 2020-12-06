FROM golang:1.14
MAINTAINER  Maksim Baldin maksim.baldin@online.liverpool.ac.uk

RUN apt-get update && apt-get install -y supervisor
RUN apt-get install -y python-pip && pip install supervisor-stdout

RUN mkdir -p /var/log/supervisor

COPY src/evaluation/cfg/service_d_supervisord.conf /etc/supervisor/conf.d/supervisord.conf

WORKDIR /service

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls

RUN go build -o service_d src/evaluation/cmd/run_service.go

CMD ["/usr/bin/supervisord"]