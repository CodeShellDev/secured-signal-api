---
title: Variables
---

# Variables

The most common type of [placeholders](../usage/advanced#placeholders) are **variables**.
These can be set under `message.variables` in your config.

> [!IMPORTANT]
> Variables are always converted into an **uppercase** string.
> Example: `number` ⇒ `NUMBER` in `{{.NUMBER}}`
> (See [Formatting](../usage/formatting#templates))

Here is an example:

```yaml
settings:
  message:
    variables:
      number: "+123400001",
      recipients: ["+123400002", "group.id", "user.id"]
```
