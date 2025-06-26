FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY dist/$TARGETOS/$TARGETARCH/app .

CMD ["./app"]