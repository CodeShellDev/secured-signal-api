FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /app/app .



FROM alpine:3.23

RUN apk --no-cache add ca-certificates

#===============================#
#        OCI Metadata           #
#===============================#

ARG IMAGE_TAG
ARG BUILD_TIME
ARG GIT_COMMIT
ARG REPOSITORY_URL

LABEL org.opencontainers.image.version=$IMAGE_TAG
LABEL org.opencontainers.image.created=$BUILD_TIME
LABEL org.opencontainers.image.revision=$GIT_COMMIT

LABEL org.opencontainers.image.source=$REPOSITORY_URL
LABEL org.opencontainers.image.url=$REPOSITORY_URL

#===============================#
#        Build Metadata         #
#===============================#

ENV IMAGE_TAG=$IMAGE_TAG
ENV BUILD_TIME=$BUILD_TIME
ENV GIT_COMMIT=$GIT_COMMIT
ENV REPOSITORY=$REPOSITORY_URL

#===============================#
#   Application Configuration   #
#===============================#

ENV DEFAULTS_PATH=/app/data/defaults.yml
ENV FAVICON_PATH=/app/data/favicon.ico

ENV CONFIG_PATH=/config/config.yml
ENV TOKENS_DIR=/config/tokens
ENV DB_PATH=/db/db.sqlite

ENV REDACT_TOKENS=true



WORKDIR /app

COPY --from=builder /app/app .
COPY data/ /app/data/

CMD ["./app"]