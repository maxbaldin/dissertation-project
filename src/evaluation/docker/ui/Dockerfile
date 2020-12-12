FROM golang:1.14
MAINTAINER  Maksim Baldin maksim.baldin@online.liverpool.ac.uk

WORKDIR /ui

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    UI_CFG=/ui/src/implementation/ui/cfg/ui.yaml


# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls

RUN go build -o ui src/implementation/ui/ui.go

CMD ["./ui"]