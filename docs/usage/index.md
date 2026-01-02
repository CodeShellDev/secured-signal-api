---
sidebar_position: 1
title: Usage
---

# Usage

In this section we'll be taking a look at how to use **Secured Signal API**.

## Basic

Here is a quick command to see if you've correctly followed the [setup instructions](../getting-started/setup):

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"message":"Hello, World!", "number":"<from>", "recipients":["<to>"]}' \
    'http://sec-signal-api:8880/v2/send'
```

This will send `Hello, World!` to `<to>` from `<from>`.

## Auth

**Secured Signal API** implements 3 auth methods:

| Method      | Example                                                    |
| :---------- | :--------------------------------------------------------- |
| Bearer Auth | Add `Authorization: Bearer API_TOKEN` to headers           |
| Basic Auth  | Add `Authorization: Basic BASE64_STRING` (`api:API_TOKEN`) |
| Query Auth  | Append `@authorization=API_TOKEN` to request URL           |

Example:

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer API_TOKEN" \
    -d '{"message":"Hello, World!", "number":"<from>", "recipients":["<to>"]}' \
    'http://sec-signal-api:8880/v2/send'
```
