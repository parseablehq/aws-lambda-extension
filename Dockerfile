FROM golang:alpine AS builder
WORKDIR /opt
COPY . /opt/
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o ./extensions/parseable-extension ./cmd/parseable-extension/main.go

FROM scratch AS base
WORKDIR /opt/extensions
COPY --from=builder /opt/extensions .
