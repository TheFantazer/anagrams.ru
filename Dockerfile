FROM golang:1.24-alpine AS builder
LABEL authors="github.com/TheFantazer"

WORKDIR /build
COPY go.mod go.su[m] ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/server .
COPY dictionaries ./dictionaries

EXPOSE 8080
CMD ["./server"]