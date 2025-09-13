FROM alpine:latest
RUN apk --no-cache add ca-certificates

ENV SERVER__PORT=8880

ENV DEFAULTS_PATH=/app/config/defaults.yml

ENV CONFIG_PATH=/config/config.yml
ENV TOKENS_DIR=/config/tokens

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

COPY dist/${TARGETOS}/${TARGETARCH}/app .

RUN ls

CMD ["./app"]
