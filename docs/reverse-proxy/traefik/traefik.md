---
title: Traefik
---

# Traefik

Want to use [**Traefik**](https://github.com/traefik/traefik) as your **reverse proxy**?
Then look no further, we'll take you through how to integrate Traefik with **Secured Signal API**.

## Prerequisites

Before moving on you must have

- already **configured** **Traefik**
- some knowledge of Traefik
- valid **SSL certificates**

## Installation

To implement Traefik in front of **Secured Signal API** you need to update your `docker-compose.yaml` file.

```yaml
services:
  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api:latest
    container_name: secured-signal
    environment:
      API__URL: http://signal-api:8080
      SETTINGS__MESSAGE__VARIABLES__RECIPIENTS: "[+123400002,+123400003,+123400004]"
      SETTINGS__MESSAGE__VARIABLES__NUMBER: "+123400001"
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

  traefik:
    image: traefik:latest
    container_name: traefik
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./config/traefik.yaml:/etc/traefik/traefik.yaml:ro
      - certs:/var/traefik/certs/:rw
      - logs:/var/log/traefik
    restart: unless-stopped
    networks:
      proxy:
        ipv4_address: 172.20.0.100

networks:
  backend:
  proxy:
```

To include the Traefik router and service labels.

## Configuration

Now you can switch over to **Secured Signal API** and add Traefik to your [trusted proxies](../../configuration/trusted-proxies.md):

```yaml
settings:
  access:
    trustedProxies:
      - 172.20.0.100
```

Then restart **Secured Signal API**:

```bash
docker compose down && docker compose up -d
```

That's it, have fun!
