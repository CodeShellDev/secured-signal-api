FROM alpine:latest
RUN apk --no-cache add ca-certificates

ENV PORT=8880

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY dist/${TARGETOS}/${TARGETARCH}/app .

RUN ls

CMD ["./app"]