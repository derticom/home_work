FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/calendar ./cmd/calendar

FROM alpine:latest
WORKDIR /root/

ENV CONFIG_PATH=./configs/calendar_config.yml

COPY --from=builder /app/calendar ./calendar

COPY configs/calendar_config.yml ./configs/calendar_config.yml
COPY migrations ./migrations/

CMD ["./calendar"]
