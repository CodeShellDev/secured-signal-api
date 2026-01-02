---
title: Field Mappings
---

# Field Mappings

A compatibility layer for **Secured Signal API**.

To improve compatibility with other services Secured Signal API provides **Field Mappings** and even a built-in `message` mapping.

<details>
<summary><strong>Default `message` mapping</strong></summary>

| Field        | Score | Field            | Score |
| ------------ | ----- | ---------------- | ----- |
| msg          | 100   | data.content     | 9     |
| content      | 99    | data.description | 8     |
| description  | 98    | data.text        | 7     |
| text         | 20    | data.summary     | 6     |
| summary      | 15    | data.details     | 5     |
| details      | 14    | body             | 2     |
| data.message | 10    | data             | 1     |

</details>

Secured Signal API will pick the highest scoring field (if available) to set the key to the correct value **using the request body**.

Field Mappings can be added by setting `message.fieldMappings` in your config:

```yaml
settings:
  message:
    fieldMappings:
      "@message":
        [
          { field: "msg", score: 80 },
          { field: "data.message", score: 79 },
          { field: "array[0].message", score: 78 },
        ]
      ".NUMBER": [{ field: "phone_number", score: 100 }]
```

> [!NOTE]
> Supported [placeholder types](../usage/advanced#placeholders):
>
> | `.` Variables | `@` Body | `#` Headers |
> | ------------- | -------- | ----------- |
> | ✅            | ✅       | ❌          |
