---
sidebar_position: 3
title: Formatting
---

# Formatting

**Secured Signal API** has some specific formatting rules to ensure for correct parsing.

## Templates

**Secured Signal API** is built with Go and therefore uses Go’s [standard templating library](https://pkg.go.dev/text/template).
As a result, any valid Go template string works in Secured Signal API.

> [!NOTE]
> The following features use Go’s templating library:
>
> - [Message Templates](../configuration/message-template)
> - [URL-to-Body Injection](./advanced#url-to-body-injection)
> - [Placeholders](./advanced#placeholders)

| Scope                                   | Example             | Note             |
| :-------------------------------------- | :------------------ | :--------------- |
| Body                                    | `{{@data.key}}`     |                  |
| Headers                                 | `{{#Content_Type}}` | `-` becomes `_`  |
| [Variables](../configuration/variables) | `{{.VAR}}`          | Always uppercase |

## String to Type

> [!TIP]
> This formatting applies to **almost every situation** where the only (allowed) **input type is a string** and **other output types are needed**

If you are using environment variables for example there would be no way of using arrays or even dictionaries as values, for these cases we have **String to Type** conversion shown below.

| Type          | Example      |
| :------------ | :----------- |
| string        | abc          |
| string        | +123         |
| int           | 123          |
| int           | -123         |
| JSON          | \{"a": "b"\} |
| array(int)    | [1, 2, 3]    |
| array(string) | [a, b, c]    |

> [!TIP]
> Escape type denotations such as `[]`, `{}`, or `-` by prefixing them with a **backslash** (`\`).  
> An **odd number** of backslashes will **escape** the character immediately following them
