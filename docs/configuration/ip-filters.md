---
title: IP Filters
---

# IP Filters

Restrict access to your **Secured Signal API** based on client IP addresses.

IP filtering allows you to explicitly **allow** or **block** requests originating from specific IPs or CIDR ranges.

## Default

By default, **all IP addresses are allowed**.

No IP-based restrictions are applied unless `access.ipFilter` is configured.

## Customize

You can modify IP access rules by configuring `access.ipFilter` in your config:

```yaml
settings:
  access:
    ipFilter:
      - "!123.456.78.9"
      - "!234.567.89.0/24"
      - 192.168.1.10
      - 10.0.0.0/24
```

By default, adding an IP or range explicitly allows it, use `!` to block it instead.

> [!IMPORTANT]
> When using `!` to block an IP or range, you must enclose it in quotes

**Supports:**

- Single IPv4 / IPv6 addresses
- CIDR notation (`10.0.0.0/24`, `2001:db8::/32`)

## Behavior

| Allow          | Block                    | Result                                    |
| -------------- | ------------------------ | ----------------------------------------- |
| `192.168.1.10` | —                        | Only `192.168.1.10` allowed               |
| —              | `!123.456.78.9`          | All allowed, except `123.456.78.9`        |
| `10.0.0.0/24`  | `!10.0.0.10`             | `10.0.0.0/24` allowed, except `10.0.0.10` |
| —              | `!0.0.0.0/0`<br/>`!::/0` | All IPv4 & IPv6 blocked                   |

### Rules

- Default: **allow all**
- Allow rules restrict access **only if no block rules exist**
- Block rules deny matching IPs
- Explicit allow overrides block
- IPv4 and IPv6 rules may be mixed

## Clients behind Proxies

For clients behind proxies IPs cannot be reliably determined without trusting the `X-Forwarded-For` HTTP header.
In order for **Secured Signal API** to trust the _XFF_ header it has to trust the request's originating proxy.

Read more about trusted proxies [here](./trusted-proxies).
