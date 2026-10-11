FROM golang:1.25.7 AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/backend .
COPY --from=builder /app/config.yaml ./

COPY web ./web

EXPOSE 8080

ENTRYPOINT ["./backend"]