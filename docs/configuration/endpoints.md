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

You can modify Blocked Endpoints by configuring `blockedEndpoints` in your config:

```yaml
settings:
  blockedEndpoints: [/v1/register, /v1/unregister, /v1/qrcodelink, /v1/contacts]
```

You can also override Blocked Endpoints by adding Allowed Endpoints to `allowedEndpoints`.

```yaml
settings:
  allowedEndpoints: [/v2/send]
```

| Config (Allow)                   | (Block)                             |   Result   |     |                   |     |
| :------------------------------- | :---------------------------------- | :--------: | --- | :---------------: | --- |
| `allowedEndpoints: ["/v2/send"]` | `unset`                             |  **all**   | 🛑  |  **`/v2/send`**   | ✅  |
| `unset`                          | `blockedEndpoints: ["/v1/receive"]` |  **all**   | ✅  | **`/v1/receive`** | 🛑  |
| `blockedEndpoints: ["/v2"]`      | `allowedEndpoints: ["/v2/send"]`    | **`/v2*`** | 🛑  |  **`/v2/send`**   | ✅  |
