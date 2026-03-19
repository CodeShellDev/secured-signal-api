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

> [!NOTE]
> The `Access-Control-Allow-Credentials` header is automatically set to `true`

**Defaults** can be set under `cors.methods` and `cors.headers`, when adding origins under `cors.origins` **overwrites** can be defined
under `methods` and `headers` like in the example above.

> [!IMPORTANT]
> The `cors.methods` setting maps to `Access-Control-Allow-Methods` and `cors.headers` to `Access-Control-Allow-Headers`
