FROM golang:1.24.5-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bot ./cmd/bot

FROM mcr.microsoft.com/playwright:v1.43.1-jammy

RUN apt-get update && apt-get install -y ca-certificates chromium-driver fonts-liberation

WORKDIR /app
COPY --from=builder /bot /bot

ENV PATH="/usr/bin/chromium-browser:${PATH}"
ENV CHROME_PATH=/usr/bin/chromium-browser

ENTRYPOINT ["/workers"]