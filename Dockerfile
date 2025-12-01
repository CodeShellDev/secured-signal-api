FROM alpine:latest
RUN apk --no-cache add ca-certificates

ARG IMAGE_TAG
ENV IMAGE_TAG=$IMAGE_TAG
LABEL org.opencontainers.image.version=$IMAGE_TAG

ENV SERVICE__PORT=8880

ENV DEFAULTS_PATH=/app/data/defaults.yml
ENV FAVICON_PATH=/app/data/favicon.ico

ENV CONFIG_PATH=/config/config.yml
ENV TOKENS_DIR=/config/tokens

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

COPY dist/${TARGETOS}/${TARGETARCH}/app .

RUN rm dist/ -r

CMD ["./app"]
