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
{{{ #://./examples/traefik.docker-compose.yaml }}}
```

To include the traefik router and service labels.

Then restart **Secured Signal API**:

```bash
docker compose down && docker compose up -d
```
