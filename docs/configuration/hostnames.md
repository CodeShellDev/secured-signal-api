---
title: Hostnames
---

# Hostnames

Hostnames can be set to create isolated realms or to restrict access by limiting to a only a small subset of hostnames.

Add hostnames, that are allowed to be used in `service.hostnames`. (default: all)

```yaml
service:
  hostnames:
    - mydomain.com
```

## Usage behind Proxy

For clients behind proxies IPs cannot be reliably determined without using the `X-Forwarded-Proto`, `X-Forwarded-Host` and `X-Forwarded-Port` HTTP headers.

For **Secured Signal API** to trust a proxy it must be added to the trusted proxies, read more [here](./trusted-proxies).
