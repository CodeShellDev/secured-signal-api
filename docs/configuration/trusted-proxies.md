---
title: Trusted Proxies
---

# Trusted Proxies

Proxies can be marked as trusted.

Add proxies to be trusted in `access.trustedProxies`.

```yaml
settings:
  access:
    trustedProxies:
      - 172.20.0.100
```

## `X-Forwarded-*` Headers

HTTP listeners only get the `proto://host:port/uri` from the incoming request, but proxies often redirect requests causing modified request URLs
`http://sec-signal-api:8880`.

To get the origin URL you have to use the `X-Forwarded-*` headers, but as you might know anyone can set headers (spoofing possible).
This means you should only trust _XF_ headers from trusted sources,
otherwise, malicious actors can change any `X-Forwarded-*` headers to be able to bypass block list, rate limits, hostname restrictions, … .

This also applies to determining the IP of a client behind a proxy, so it is extremely important to allow for using the _XF_ headers when a proxy is trusted.
