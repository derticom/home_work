FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/sender ./cmd/sender

FROM alpine:latest
WORKDIR /root/

ENV CONFIG_PATH=./configs/sender_config.yml

COPY --from=builder /app/sender ./sender
COPY configs/sender_config.yml ./configs/sender_config.yml

CMD ["./sender"]
