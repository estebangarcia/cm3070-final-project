FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o registry main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/registry .

ENTRYPOINT ["./registry"]