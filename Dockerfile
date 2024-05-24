FROM golang:1.22.1 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app cmd/booker/main.go

FROM alpine:latest
WORKDIR /opt/booker
COPY --from=builder /build/app app
