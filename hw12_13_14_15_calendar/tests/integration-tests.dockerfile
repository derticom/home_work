FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ .

RUN CGO_ENABLED=0 GOOS=linux go test -c -o integration-tests ./tests

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/integration-tests ./integration-tests

CMD ["./integration-tests", "-test.v"]
