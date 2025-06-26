FROM golang:1.24 AS builder

WORKDIR /build
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o app .

FROM alpine:latest AS bu

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/app .

CMD ["./app"]