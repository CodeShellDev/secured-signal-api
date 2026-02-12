---
sidebar_position: 2
title: Advanced
---

# Advanced

Here you will be explained all the neat tricks and quirks for **Secured Signal API**

## Placeholders

Placeholders do exactly what you think they do: They **replace** actual values.
These can be especially **helpful** when managing **variables** across multiple services.

### How to Use

| Scope                                  | Example                                                                      |
| :------------------------------------- | :--------------------------------------------------------------------------- |
| Body                                   | **Short:** `{{@data.key}}`<br/>**Standard:** `{{.body.data.key}}`            |
| Header                                 | **Short:** `{{#Content_Type}}`<br/>**Standard:** `{{.headers.Content_Type}}` |
| [Variable](../configuration/variables) | **`{{.VAR}}`**                                                               |

> [!NOTE]
> Formatting rules (capitalization, escaping, and typing) are defined in [Formatting](./formatting)

> [!TIP]
> Sometimes using the **Stable** placeholder mode is better, for example: `{{ index . ".body.data.key" }}`

### Where to Use

| Scope | Example                                                          |
| :---- | :--------------------------------------------------------------- |
| Body  | `{"number": "{{ .NUMBER }}", "recipients": "{{ .RECIPIENTS }}"}` |
| Query | `http://sec-signal-api:8880/v1/receive/?@number={{.NUMBER}}`     |
| Path  | `http://sec-signal-api:8880/v1/receive/{{.NUMBER}}`              |

**Combine them:**

```json
"message": "{{.NUMBER}} -> {{.RECIPIENTS}}"
```

**Mix and match:**

```json
"message": "{{#X_Forwarded_For}} just send from {{.NUMBER}}"
```

> [!NOTE]
> Placeholders follow strict formatting rules ([See Formatting](./formatting#templates))

## URL-to-Body Injection

> _Sounds scary… but it really isn't._ 🫣

**URL-to-Body Injection** is a powerful feature designed for **restricted or inflexible environments**.

In some setups, webhook configuration is extremely limited.
For example, you may **only** be able to define a webhook URL, without any control over the **request body**.
This becomes a problem when every receiving service is expected to support the **Signal CLI REST API** format.
In such cases, using a simple, generic webhook is not possible.

**URL-to-Body Injection** solves this by allowing you to inject values directly into the request body via query parameters or seperated path parameters.

`http://sec-signal-api:8880/@key2=value2/?@key=value`

> [!IMPORTANT]
>
> - Supported value types include **strings**, **integers**, **arrays**, and **JSON objects**
> - See [Formatting](./formatting#string-to-type) for details on supported structures and syntax

> [!NOTE]
> Supported [placeholder types](./advanced#placeholders):
>
> | `.` Variables | `@` Body | `#` Headers |
> | ------------- | -------- | ----------- |
> | ❌            | ✅       | ❌          |
