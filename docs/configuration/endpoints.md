---
title: Endpoints
---

# Endpoints

Restrict access to your **Secured Signal API**.

## Default

Secured Signal API is just a Proxy, which means any and **all** of the **Signal CLI REST API** **endpoints are available**,
but by default the following endpoints are **blocked**, because of Security Concerns:

| Endpoint              |                    |
| :-------------------- | ------------------ |
| **/v1/about**         | **/v1/unregister** |
| **/v1/configuration** | **/v1/qrcodelink** |
| **/v1/devices**       | **/v1/contacts**   |
| **/v1/register**      | **/v1/accounts**   |

## Customize

> [!NOTE]
> Matching works by checking if the requested Endpoints starts with a Blocked or an Allowed Endpoint

You can modify endpoints by configuring `access.endpoints` in your config:

```yaml
settings:
  access:
    endpoints:
      - !/v1/register
      - !/v1/unregister
      - !/v1/qrcodelink
      - !/v1/contacts
      - /v2/send
```

By default adding an endpoint explictly allows access to it, use `!` to block it instead.

| Config (Allow) | (Block)        |   Result   |     |                   |     |
| :------------- | :------------- | :--------: | --- | :---------------: | --- |
| `/v2/send`     | `unset`        |  **all**   | 🛑  |  **`/v2/send`**   | ✅  |
| `unset`        | `!/v1/receive` |  **all**   | ✅  | **`/v1/receive`** | 🛑  |
| `/v2`          | `!/v2/send`    | **`/v2*`** | 🛑  |  **`/v2/send`**   | ✅  |
