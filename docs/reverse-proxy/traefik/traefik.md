---
title: Traefik
---

# Traefik

Want to use [**traefik**](https://github.com/traefik/traefik) as your **Reverse Proxy**?
Then look no further, we'll take you through how to integrate traefik with **Secured Signal API**.

## Prerequisites

Before moving on you must have

- already **configured** **traefik**
- some knowledge of traefik
- valid **SSL Certificates**

## Installation

To implement traefik infront of **Secured Signal API** you need to update your `docker-compose.yaml` file.

```yaml
services:
  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api:latest
    container_name: secured-signal
    environment:
      API__URL: http://signal-api:8080
      SETTINGS__VARIABLES__RECIPIENTS: "[+123400002,+123400003,+123400004]"
      SETTINGS__VARIABLES__NUMBER: "+123400001"
      API__TOKENS: "[LOOOOOONG_STRING]"
    labels:
      - traefik.enable=true
      - traefik.http.routers.signal-api.rule=Host(`signal-api.mydomain.com`)
      - traefik.http.routers.signal-api.entrypoints=websecure
      - traefik.http.routers.signal-api.tls=true
      - traefik.http.routers.signal-api.tls.certresolver=cloudflare
      - traefik.http.routers.signal-api.service=signal-api-svc
      - traefik.http.services.signal-api-svc.loadbalancer.server.port=8880
      - traefik.docker.network=proxy
    restart: unless-stopped
    networks:
      proxy:
      backend:
        aliases:
          - secured-signal-api

networks:
  backend: {}
  proxy:
    external: true
```

To include the traefik router and service labels.

Then restart **Secured Signal API**:

```bash
docker compose down && docker compose up -d
```
