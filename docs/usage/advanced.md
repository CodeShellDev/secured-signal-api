---
sidebar_position: 2
title: Advanced
---

# Advanced

Here you will be explained all of the neat tricks and quirks for **Secured Signal API**

## Placeholders

Placeholders do exactly what you think they do: They **replace** actual values.
These can be especially **helpful** if have to manage multiple **Variables** and don't want to hardcode them into your request every time.

### How to use

| Type                   | Example             | Note             |
| :--------------------- | :------------------ | :--------------- |
| Body                   | `{{@data.key}}`     |                  |
| Header                 | `{{#Content_Type}}` | `-` becomes `_`  |
| [Variable](#variables) | `{{.VAR}}`          | always uppercase |

### Where to use

| Type  | Example                                                          |
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

## KeyValue Pair Injection

> _OoOhhh scary..._ 🫣

They may sound a bit dangerous (and can be), but **KeyValue Pair Injections** are **extremely useful** in **limited environments**.

In some setups you could be dealing with a very **limited control of webhooks**, for example you only can set a webhook url and cannot modify the body.
This is very annoying since this means **every programm** needs to support **Signal CLI REST API** and you cannot just use a **generic webhook**.
This is why we have **KeyValue Pair Injection** which lets you inject query values into the request's body:

`http://sec-signal-api:8880/?@key=value`

> [!IMPORTANT]
> Prefix with `@` for injecting into the Body.
> Supported types include **strings**, **ints**, **arrays** and **json dictionaries**. See [Formatting](./formatting).
