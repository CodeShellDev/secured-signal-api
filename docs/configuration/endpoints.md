---
title: Endpoints
---

# Endpoints

Restrict access to your **Secured Signal API**.

## Default

Secured Signal API is just a proxy, which means any and all the **Signal CLI REST API** **endpoints are available**,
because of security concerns the following endpoints are blocked:

| Endpoint              |                    |
| :-------------------- | ------------------ |
| **/v1/configuration** | **/v1/unregister** |
| **/v1/devices**       | **/v1/contacts**   |
| **/v1/register**      | **/v1/accounts**   |
| **/v1/qrcodelink**    |                    |

## Customize

> [!IMPORTANT]
>
> 1. Matching uses [regex](https://regex101.com)
> 2. On error [glob-style patterns](https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html) are used instead

> [!NOTE]
> Quick reminder, how [glob-style patterns](https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html) work:
>
> - `*` matches any sequence of characters
> - `?` matches a single character
> - `[abc]` matches one of the characters in the brackets

You can modify endpoints by configuring `access.endpoints` in your config:

```yaml
settings:
  access:
    endpoints:
      - "!/v1/register"
      - "!/v1/unregister"
      - "!/v1/qrcodelink"
      - "!/v1/contacts"
      - /v2/send
```

By default, adding an endpoint explicitly allows access to it, use `!` to block it instead.

> [!IMPORTANT]
> When using `!` to block you must enclose the endpoint in quotes, like in the example above

## Behavior

| Allow      | Block          | Result                                    |
| ---------- | -------------- | ----------------------------------------- |
| `/v2/send` | —              | **Only** `/v2/send` allowed               |
| —          | `!/v1/receive` | **All** allowed, **except** `/v1/receive` |
| `/v2/send` | `!/v2/.*`      | **Only** `/v2/send` allowed               |

### Rules

- Default: **allow all**
- Allow rules exist: default **block**
- Only block rules exist: default **allow**
- Explicit allow **overrides** block
