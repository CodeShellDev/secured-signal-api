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
        [
          "Content-Type",
          "Content-Language",
          "Authorization",
          "Accept",
          "Accept-Language",
        ]
        origins:
          - url: "https://domain.com"
          - url: "https://example.com/path"
            methods: [GET]
            headers: ["Content-Type"]
```

The `cors.methods` and `cors.headers` settings act as **defaults** for origins, which do not **overwrite** `methods` or `headers`.

> [!NOTE]
> Defaults for `cors.methods` and `cors.headers` are **already** defined as in the above

> [!IMPORTANT]
> During preflight requests (`OPTIONS`) no authentication can be provided, this means using [**token configs**](../#token-configs) is **not possible**, use the [main config](../#config-files) instead
