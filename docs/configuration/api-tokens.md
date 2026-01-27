---
title: API Tokens
---

# API Tokens

> [!IMPORTANT]
> Using API tokens is highly **recommended**, but not mandatory.
> Some important Security Features won't be available (for example the default [blocked endpoints](./endpoints#default))

Define API tokens for accessing **Secured Signal API**.

```yaml
api:
  tokens: [token1, token2, token3]
```

See [Auth](./auth) for possible auth methods and how to activate additional ones.

> [!NOTE]
> Blocked endpoints can be reactivated by manually configuring them
