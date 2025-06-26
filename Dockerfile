FROM alpine:latest
RUN apk --no-cache add ca-certificates

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY dist/${TARGETOS}/${TARGETARCH}/app .

RUN ls

CMD ["./app"]