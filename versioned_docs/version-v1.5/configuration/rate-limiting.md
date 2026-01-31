---
title: Rate Limiting
---

# Rate Limiting

Rate limiting is used to control how many requests or actions a token or client can perform within a given period.
This helps prevent abuse, protect downstream services, and ensure fair usage.

Rate limits can be defined in `settings.access.rateLimiting`:

```yaml
settings:
  access:
    rateLimiting:
      limit: 100
      period: 1m
```

- `limit`: The maximum number of allowed requests in the given period
- `period`: The duration over which the limit is measured (supports Go duration format like `1m`, `10s`, `1h`)

When a request exceeds the configured rate limit the server responds with `429` `Too Many Requests`.

> [!NOTE]
>
> [Trusted clients](./trusted-ips.md) are allowed to bypass rate limits
