FROM golang:1.14
MAINTAINER  Maksim Baldin maksim.baldin@online.liverpool.ac.uk

RUN apt-get update && apt-get install -y supervisor
RUN apt-get install -y python-pip && pip install supervisor-stdout
RUN apt-get install libpcap-dev -y

RUN mkdir -p /var/log/supervisor

COPY src/evaluation/cfg/service_a_supervisord.conf /etc/supervisor/conf.d/supervisord.conf

WORKDIR /service

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls

RUN go build -o service_a src/evaluation/cmd/run_service.go
RUN go build -o agent src/implementation/agent/agent.go

CMD ["/usr/bin/supervisord"]