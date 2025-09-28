## Getting Started

> **Prerequisites**: You need Docker and Docker Compose installed.

Get the latest version of the `docker-compose.yaml` file:

```yaml
services:
  signal-api:
    image: bbernhard/signal-cli-rest-api:latest
    container_name: signal-api
    environment:
      - MODE=normal
    volumes:
      - ./data:/home/.local/share/signal-cli
    restart: unless-stopped
    networks:
      backend:
        aliases:
          - signal-api

  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api:latest
    container_name: secured-signal
    environment:
      API__URL: http://signal-api:8080
      SETTINGS__VARIABLES__RECIPIENTS:
        '[+123400002, +123400003, +123400004]'
      SETTINGS__VARIABLES__NUMBER: "+123400001"
      API__TOKENS: '[LOOOOOONG_STRING]'
    ports:
      - "8880:8880"
    restart: unless-stopped
    networks:
      backend:
        aliases:
          - secured-signal-api

networks:
  backend:
```

And add secure Token(s) to `api.tokens`. See [API Tokens](#api-tokens).

> [!IMPORTANT]
> In this documentation, we use `sec-signal-api:8880` as the host for simplicity.
> Replace it with your actual container/host IP, port, or hostname.
> 