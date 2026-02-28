---
title: Injecting
---

# Injecting

Configure injection settings.

**Example:**

```yaml
settings:
  message:
    injecting:
      urlToBody:
        query: true
        path: true
```

The example config entry above enables [injection](../usage/advanced#url-to-body-injection) from the following mediums: from the **query** and also from **path**:

```bash
curl -H "Authorization: Bearer API_TOKEN" \
    'http://sec-signal-api:8880/@number={{.NUMBER}}/v1/send?@message={{.MESSAGE}}'
```
