---
title: Caddy
---

# Caddy

Want to use [**Caddy**](https://github.com/caddyserver/caddy) as your **reverse proxy**?
These instructions will take you through the steps.

## Prerequisites

Before moving on you must have

- some knowledge of **Caddy**
- already deployed **Caddy**

## Installation

Add Caddy to your `docker-compose.yaml` file.

```yaml
{{{ #://./examples/caddy.docker-compose.yaml }}}
```

## Setup

Create a `Caddyfile` in your `docker-compose.yaml` folder and mount it to `/etc/caddy/Caddyfile` in your Caddy container.

```apacheconf
{{{ #://./examples/Caddyfile }}}
```

## Configuration

Now you can switch over to **Secured Signal API** and add Caddy to your [trusted proxies](../../configuration/trusted-proxies.md):

```yaml
settings:
  access:
    trustedProxies:
      - 172.20.0.100
```

Then spin up your stack:

```bash
docker compose up -d
```

And you are ready to go!
