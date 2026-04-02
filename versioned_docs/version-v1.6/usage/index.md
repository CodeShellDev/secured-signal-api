---
sidebar_position: 1
title: Usage
---

# Usage

In this section we'll be taking a look at how to use **Secured Signal API**.

## Basic

Here is a quick command to see if you've correctly followed the [setup instructions](./getting-started/setup):

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"message":"Hello, World!", "number":"<from>", "recipients":["<to>"]}' \
    'http://sec-signal-api:8880/v2/send'
```

This will send `Hello, World!` to `<to>` from `<from>`.

## Auth

**Secured Signal API** implements 5 auth methods:

| Method      | Example                                                                                                      |
| :---------- | :----------------------------------------------------------------------------------------------------------- |
| Bearer Auth | `Authorization: Bearer API_TOKEN` (header)                                                                   |
| Basic Auth  | `Authorization: Basic base64(api:API_TOKEN)` (header)<br/>`http://api:API_TOKEN@host:port` (client specific) |
| Query Auth  | `http://host:port/abc?@auth=API_TOKEN` (query parameter)                                                     |
| Path Auth   | `http://host:port/@auth=API_TOKEN/abc` (path parameter)                                                      |
| Body Auth   | `{ "auth": "API_TOKEN" }` (request body field)                                                               |

> [!WARNING]
> **Query** and **Path** auth are disabled by default and [must be enabled in the config](../configuration/auth)

**Example:**

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer API_TOKEN" \
    -d '{"message":"Hello, World!", "number":"<from>", "recipients":["<to>"]}' \
    'http://sec-signal-api:8880/v2/send'
```
