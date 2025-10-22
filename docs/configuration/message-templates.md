---
title: Message Templates
---

# Message Templates

Message Templates are the best way to **structure** and **customize** your messages and also very useful for **compatiblity** between different services.

Configure them by using the `messageTemplates` attribute in you config.

These support Go Templates (See [Usage](../usage/advanced)) and work by templating the `message` attribute in the request's body.

Here is an example:

```yaml
{{{ #://./examples/message-template.yml }}}
```

> [!IMPORTANT]
> Message Templates support [Standard Golang Templating](../usage/advanced).
> Use `@data.key` to reference Body Keys, `#Content_Type` for Headers and `.KEY` for [Variables](./variables).
