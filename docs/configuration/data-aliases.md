---
title: Data Aliases
---

# Data Aliases

A Compatibility Layer for **Secured Signal API**.

To improve compatibility with other services Secured Signal API provides **Data Aliases** and even a built-in `message` Alias.

<details>
<summary><strong>Default `message` Aliases</strong></summary>

| Alias        | Score | Alias            | Score |
| ------------ | ----- | ---------------- | ----- |
| msg          | 100   | data.content     | 9     |
| content      | 99    | data.description | 8     |
| description  | 98    | data.text        | 7     |
| text         | 20    | data.summary     | 6     |
| summary      | 15    | data.details     | 5     |
| details      | 14    | body             | 2     |
| data.message | 10    | data             | 1     |

</details>

Secured Signal API will pick the highest scoring **Data Alias** (if available) to set the key to the correct value **using the request body**.

Data Aliases can be added by setting `dataAliases` in your config:

```yaml
settings:
  dataAliases:
    "@message":
      [
        { alias: "msg", score: 80 },
        { alias: "data.message", score: 79 },
        { alias: "array[0].message", score: 78 },
      ]
    ".NUMBER": [{ alias: "phone_number", score: 100 }]
```

> [!IMPORTANT]
> Use `@` for aliasing Body Keys and `.` for aliasing Variables.
