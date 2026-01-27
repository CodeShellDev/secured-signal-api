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
{{{ #://./examples/traefik.docker-compose.yaml }}}
```

To include the Traefik router and service labels.

Then restart **Secured Signal API**:

```bash
docker compose down && docker compose up -d
```
