---
sidebar_position: 3
title: Formatting
---

# Formatting

**Secured Signal API** has some specific formatting rules to ensure for correct parsing.

## Templates

**Secured Signal API** is built with Go and hence uses Go's [Standard Templating Library](https://pkg.go.dev/text/template).
Which means that any valid Go template string will also work in Secured Signal API.

> [!NOTE]
> Go's templating library is used in the following features:
> <br/>- [Message Templates](../configuration/message-template) <br/>- [Placeholders](./advanced)

But you will mostly be using `{{.VAR}}`.

## String to Type

> [!TIP]
> This formatting applies to **almost every situation** where the only (allowed) **Input Type** is a string and **other** **Output Types** are **needed**.

If you are using environment variables for example there would be no way of using arrays or even dictionaries as values, for these cases we have **String to Type** conversion shown below.

| type       | example             |
| :--------- | :------------------ |
| string     | abc                 |
| string     | +123                |
| int        | 123                 |
| int        | -123                |
| json       | \{"a":"b","c":"d"\} |
| array(int) | [1,2,3]             |
| array(str) | [a,b,c]             |

> [!NOTE]
> Escape Type Denotations, like `[]` or `{}` (also `-`) with a `\` **backslash**.
> An **odd** number of **backslashes** **escape** the character in front of them.
