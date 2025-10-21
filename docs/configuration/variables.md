---
title: Variables
---

# Variables

The most common type of [Placeholders](../usage/advanced) are **variables**.
Which can be set under `variables` in your config.

> [!IMPORTANT]
> Every Placeholder Key will be converted into an **uppercase** string.
> Example: `number` becomes `NUMBER` in `{{.NUMBER}}`

Here is an example:

```yaml
settings:
  variables:
    number: "+123400001",
    recipients: ["+123400002", "group.id", "user.id"]
```
