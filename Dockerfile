# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && \
    apt-get install -y libportaudio2 portaudio19-dev

# Build all commands
RUN CGO_ENABLED=1 go build -o bin/collection ./cmd/collection/main.go && \
    go build -o bin/client ./cmd/client/main.go && \
    go build -o bin/effect ./cmd/effect/main.go && \
    go build -o bin/server ./cmd/server/main.go && \
    go build -o bin/user ./cmd/user/main.go

# Final stage
FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    libportaudio2 \
    portaudio19-dev \
    pkg-config \
    git

COPY --from=builder /app/bin ./bin
COPY .env ./

CMD ["sh", "-c", "./bin/server & ./bin/client && wait"]
