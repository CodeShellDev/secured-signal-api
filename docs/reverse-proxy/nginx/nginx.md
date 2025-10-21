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
file not found: /home/runner/work/secured-signal-api/secured-signal-api/examples/nginx.docker-compose.yaml
```

To include the needed mounts for your certificates and your config.

Create a `nginx.conf` file in the `docker-compose.yaml` folder and mount it to `/etc/nginx/conf.d/default.conf` in your nginx container.

```conf
file not found: /home/runner/work/secured-signal-api/secured-signal-api/examples/nginx.conf
```

Add your `cert.key` and `cert.crt` into your `certs/` folder and mount it to `/etc/nginx/ssl`.

Lastly spin up your stack:

```bash
docker compose up -d
```

And you are ready to go!
