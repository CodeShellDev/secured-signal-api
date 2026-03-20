---
title: CORS
---

# CORS

Configure CORS headers via the `settings.access.cors` setting.

**Example:**

```yaml
settings:
  access:
    cors:
      methods: [GET, POST, PUT, PATCH, DELETE, OPTIONS]
      headers:
        ["Content-Type", "Authorization", "Accept", "Accept-Language", "Origin"]
        origins:
          - url: "https://domain.com"
          - url: "https://example.com/path"
            methods: [GET]
            headers: ["Content-Type"]
```

The `cors.methods` and `cors.headers` settings act as **defaults** for origins, that do not **overwrite** `methods` or `headers`.

> [!NOTE]
> Defaults for `cors.methods` and `cors.headers` are **already** defined as in the above
