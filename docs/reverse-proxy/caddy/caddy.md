---
title: Caddy
---

# Caddy

Want to use [**caddy**](https://github.com/caddyserver/caddy) as your **Reverse Proxy**?
These instructions will take you through the steps.

## Prerequisites

Before moving on you must have

- some knowledge of **caddy**
- already deployed **caddy**

## Installation

Add caddy to your `docker-compose.yaml` file.

```yaml
services:
  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api:latest
    container_name: secured-signal-api
    environment:
      API__URL: http://signal-api:8080
      SETTINGS__MESSAGE__VARIABLES__RECIPIENTS: "[+123400002,+123400003,+123400004]"
      SETTINGS__MESSAGE__VARIABLES__NUMBER: "+123400001"
      API__TOKENS: "[LOOOOOONG_STRING]"
    restart: unless-stopped
    networks:
      backend:
        aliases:
          - secured-signal-api

  caddy:
    image: caddy:latest
    container_name: caddy
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - data:/data
    depends_on:
      - secured-signal

networks:
  backend: {}

volumes:
  data: {}
```

Create a `Caddyfile` in your `docker-compose.yaml` folder and mount it to `/etc/caddy/Caddyfile` in your caddy container.

```apacheconf
# Replace with your actual domain
domain.com {
    # Use whatever network alias you set in the docker-compose file
    reverse_proxy secured-signal-api:8880

    # Optional: basic security headers
    header {
        Strict-Transport-Security "max-age=31536000;"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "DENY"
        Referrer-Policy "no-referrer"
    }
}

# HTTP redirect to HTTPS
http://domain.com {
    redir https://{host}{uri} permanent
}
```

Then spin up your stack:

```bash
docker compose up -d
```

And you are ready to go!
