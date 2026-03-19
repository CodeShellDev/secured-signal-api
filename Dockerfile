FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o app .

FROM alpine:3.22

RUN apk --no-cache add ca-certificates

ARG IMAGE_TAG
ENV IMAGE_TAG=$IMAGE_TAG
LABEL org.opencontainers.image.version=$IMAGE_TAG

ENV SERVICE__PORT=8880

ENV DEFAULTS_PATH=/app/data/defaults.yml
ENV FAVICON_PATH=/app/data/favicon.ico

ENV CONFIG_PATH=/config/config.yml
ENV TOKENS_DIR=/config/tokens

ENV DB_PATH=/db/db.sqlite3

ENV REDACT_TOKENS=true

WORKDIR /app

COPY --from=builder /app/app .

COPY data/ /app/data/

CMD ["./app"]
