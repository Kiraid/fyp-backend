# Stage 1: Build
FROM golang:1.22.3 AS builder

RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev

WORKDIR /app
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main -ldflags="-s -w"

# Stage 2: Minimal runtime
FROM ubuntu:latest
WORKDIR /root/
COPY --from=builder /app/main .
RUN apt-get update && apt-get install -y ca-certificates
RUN chmod +x main


CMD ["./main"]
