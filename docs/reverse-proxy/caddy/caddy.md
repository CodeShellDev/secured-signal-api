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

Create or update your `Caddyfile` file and mount it to `/etc/caddy/Caddyfile` in your caddy container.

```conf
{{{ #://./examples/Caddyfile }}}
```

Then spin up your stack:

```bash
docker compose up -d
```
