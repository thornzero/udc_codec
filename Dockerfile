# Stage 1: Build
FROM golang:1.24.3 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN make build

# Stage 2: Runtime
FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/bin/server ./server
COPY .env .env
COPY data/ ./data/
COPY tags.db ./tags.db
COPY --from=builder /app/webserver ./webserver
COPY ./frontend ./frontend

EXPOSE 8080

CMD ["./webserver"]