FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/scheduler ./cmd/scheduler

FROM alpine:latest
WORKDIR /root/

ENV CONFIG_PATH=./configs/scheduler_config.yml

COPY --from=builder /app/scheduler ./scheduler
COPY configs/scheduler_config.yml ./configs/scheduler_config.yml

CMD ["./scheduler"]
