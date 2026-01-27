---
title: Message Template
---

# Message Template

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
