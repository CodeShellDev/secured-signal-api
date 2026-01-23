---
title: Auth
---

# Auth

The `auth` setting under `api` is used for configuring auth methods and setting API tokens.
(default methods: `bearer, basic, body`)

**Example:**

```yaml
api:
  auth:
    methods: [query, path]
    tokens:
      - set: [token4, token5]
        methods: [body]
```

In the example above **Query** and **Path** auth have been enabled in favor of **Bearer**, **Basic** and **Body** auth.
This applies to any token defined in [`api.tokens`](./api-tokens), <br/>`token4` and `token5` on the other hand only allow for being used with **Body** auth.
