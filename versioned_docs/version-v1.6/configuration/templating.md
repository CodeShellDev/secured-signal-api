---
title: Templating
---

# Templating

Configure templating settings.

**Example:**

```yaml
settings:
  message:
    templating:
      body: true
      query: true
      path: true
```

The example config entry above enables [templating or rather placeholders](../usage/advanced#placeholders) in all mediums: in the **body**, in the **query** and also in the **path**:

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer API_TOKEN" \
    -d '{ "number": "{{.NUMBER}}" }' \
    'http://sec-signal-api:8880/{{.PATH}}/v1/send?@message={{.MESSAGE}}'
```
