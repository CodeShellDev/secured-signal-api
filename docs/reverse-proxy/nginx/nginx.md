---
title: NGINX
---

# NGINX

Want to use [**Nginx**](https://github.com/nginx/nginx) as your **reverse proxy**?
No problem here are the instructions…

## Prerequisites

Before moving on you must have

- some knowledge of **Nginx**
- valid **SSL certificates**

## Installation

To implement Nginx in front of **Secured Signal API** you need to update your `docker-compose.yaml` file.

```yaml
{{{ #://./examples/nginx.docker-compose.yaml }}}
```

To include the needed mounts for your certificates and your config.

## Setup

Create a `nginx.conf` file in the `docker-compose.yaml` folder and mount it to `/etc/nginx/conf.d/default.conf` in your Nginx container.

```nginx
{{{ #://./examples/nginx.conf }}}
```

Add your `cert.key` and `cert.crt` into your `certs/` folder and mount it to `/etc/nginx/ssl`.

## Configuration

Now you can switch over to **Secured Signal API** and add Nginx to your [trusted proxies](../../configuration/trusted-proxies.md):

```yaml
settings:
  access:
    trustedProxies:
      - 172.20.0.100
```

Lastly spin up your stack:

```bash
docker compose up -d
```

And you are ready to go!
