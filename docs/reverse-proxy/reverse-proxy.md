### Reverse Proxy

#### Traefik

Take a look at the [traefik](https://github.com/traefik/traefik) implementation:

```yaml
services:
  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api:latest
    container_name: secured-signal
    environment:
      API__URL: http://signal-api:8080
      SETTINGS__VARIABLES__RECIPIENTS:
        '[+123400002,+123400003,+123400004]'
      SETTINGS__VARIABLES__NUMBER: "+123400001"
      API__TOKENS: '[LOOOOOONG_STRING]'
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
  backend:
  proxy:
    external: true
```

#### NGINX Proxy

This is the [NGINX](https://github.com/nginx/nginx) `docker-compose.yaml` file:

```yaml
services:
  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api:latest
    container_name: secured-signal-api
    environment:
      API__URL: http://signal-api:8080
      SETTINGS__VARIABLES__RECIPIENTS: "[+123400002,+123400003,+123400004]"
      SETTINGS__VARIABLES__NUMBER: "+123400001"
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
    restart: unless-stopped
    networks:
      frontend:
      backend:

networks:
  backend:
  frontend:
```

Create a `nginx.conf` file in the `docker-compose.yaml` folder and mount it to `etc/nginx/conf.d/default.conf`:

```conf
server {
    # Allow SSL on Port 443
    listen 443 ssl;

    # Add allowed hostnames which nginx should respond to
    # `_` for any
    server_name localhost;

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
    server_name localhost;
    return 301 https://$host$request_uri;
}
```

Lastly add your `cert.key` and `cert.crt` into your `certs/` folder and mount it to `/etc/nginx/ssl`.