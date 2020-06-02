FROM golang:1.13 as builder

WORKDIR /app

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY pkg/ ./pkg

FROM alpine:latest