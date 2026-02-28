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

## Message Templates

**Message Templates** are the best way to **structure** and **customize** your messages and can be very useful for **compatibility** between different services.

Configure them by using the `message.template` attribute in you config.

These support Go templates (see [Formatting](../usage/formatting#templates)) and work by templating the `message` field in the request's body.

Here is an example:

```yaml
{{{ #://./examples/message-template.yml }}}
```

> [!NOTE]
> Supported [placeholder types](../usage/advanced#placeholders):
>
> | `.` Variables | `@` Body | `#` Headers |
> | ------------- | -------- | ----------- |
> | ✅            | ✅       | ✅          |
