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

You can modify endpoints by configuring `access.endpoints` in your config:

```yaml
settings:
  access:
    endpoints:
      blocked:
        - pattern: /v1/register
          matchType: prefix
        - pattern: /v1/unregister
          matchType: prefix
        - pattern: /v1/qrcodelink
          matchType: prefix
        - pattern: /v1/contacts
          matchType: prefix
      allowed:
        - /v2/send
```

## Match Types

Available options for `matchType` are:

{{{ #://../templates/match-rules.string.md.tmpl }}}

## Behavior

| Allow      | Block              | Result                                    |
| ---------- | ------------------ | ----------------------------------------- |
| `/v2/send` | —                  | **Only** `/v2/send` allowed               |
| —          | `/v1/receive`      | **All** allowed, **except** `/v1/receive` |
| `/v2/send` | `/v2/.*` (`regex`) | **Only** `/v2/send` allowed               |

### Rules

{{{ #://../templates/block-rules.md.tmpl }}}
