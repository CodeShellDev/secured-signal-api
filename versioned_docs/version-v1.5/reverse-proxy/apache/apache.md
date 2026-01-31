---
title: Apache
---

# Apache

Want to use [**Apache**](https://github.com/apache/apache) as your **reverse proxy**?
Then your in luck you've come to the right place!

## Prerequisites

Before moving on you must have

- some knowledge of **Apache**
- valid **SSL certificates**

## Installation

To implement Apache in front of **Secured Signal API** you need to update your `docker-compose.yaml` file.

```yaml
{{{ #://./examples/apache.docker-compose.yaml }}}
```

## Setup

Create a `apache.conf` file in the `docker-compose.yaml` folder and mount it to `/usr/local/apache2/conf.d` in your Apache container.

```apacheconf
{{{ #://./examples/apache.conf }}}
```

Add your `cert.key` and `cert.crt` into your `certs/` folder and mount it to `/etc/ssl`.

## Configuration

Now you can switch over to **Secured Signal API** and add Apache to your [trusted proxies](../../configuration/trusted-proxies.md):

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
