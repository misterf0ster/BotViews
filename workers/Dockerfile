FROM golang:1.24.5-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/main.go

FROM mcr.microsoft.com/playwright:v1.52.0-jammy

RUN useradd -m -s /bin/bash botuser

WORKDIR /app

COPY --from=builder /app/bot .

RUN chmod +x ./bot

USER botuser

ENTRYPOINT ["./bot"]
