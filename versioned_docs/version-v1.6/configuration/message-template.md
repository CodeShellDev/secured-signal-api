---
title: Message Template
---

# Message Template

**Message Template** are the best way to **structure** and **customize** your messages and can be very useful for **compatibility** between different services.

Configure them by using the `message.messageTemplate` setting in you config.

These support Go templates (see [Formatting](../usage/formatting#templates)) and work by templating the `message` field in the request's body.

Here is an example:

```yaml
+{{{ read "./examples/message-template.yml" }}}
```

+{{{ readArgs "../templates/request-keys.md.gtmpl" "variables" "body" "headers" }}}
