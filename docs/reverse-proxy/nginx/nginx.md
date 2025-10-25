---
title: NGINX
---

# NGINX

Want to use [**nginx**](https://github.com/nginx/nginx) as your **Reverse Proxy**?
No problem here are the instructions.

## Prerequisites

Before moving on you must have

- some knowledge of **nginx**
- valid **SSL Certificates**

## Installation

To implement nginx infront of **Secured Signal API** you need to update your `docker-compose.yaml` file.

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

  nginx:
    image: nginx:latest
    container_name: secured-signal-proxy
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf
      # Load SSL certificates: cert.key, cert.crt
      - ./certs:/etc/nginx/ssl
    ports:
      - "443:443"
      - "80:80"
    depends_on:
      - secured-signal
    restart: unless-stopped
    networks:
      backend:

networks:
  backend: {}
```

To include the needed mounts for your certificates and your config.

Create a `nginx.conf` file in the `docker-compose.yaml` folder and mount it to `/etc/nginx/conf.d/default.conf` in your nginx container.

```apacheconf
server {
    # Allow SSL on Port 443
    listen 443 ssl;

    # Add allowed hostnames which nginx should respond to
    # `_` for any
    server_name domain.com;

    ssl_certificate /etc/nginx/ssl/cert.crt;
    ssl_certificate_key /etc/nginx/ssl/cert.key;

    location / {
        # Use whatever network alias you set in the docker-compose file
        proxy_pass http://secured-signal-api:8880;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Fowarded-Proto $scheme;
    }
}

# Redirect HTTP to HTTPs
server {
    listen 80;
    server_name domain.com;
    return 301 https://$host$request_uri;
}
```

Add your `cert.key` and `cert.crt` into your `certs/` folder and mount it to `/etc/nginx/ssl`.

Lastly spin up your stack:

```bash
docker compose up -d
```

And you are ready to go!
