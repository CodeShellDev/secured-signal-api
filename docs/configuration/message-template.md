---
title: Message Template
---

# Message Template

**Message Templates** are the best way to **structure** and **customize** your messages and can be very useful for **compatibility** between different services.

Configure them by using the `message.template` attribute in you config.

These support Go templates (see [Formatting](../usage/formatting#templates)) and work by templating the `message` field in the request's body.

Here is an example:

```yaml
settings:
  message:
    template: |
      {{- $greeting := "Hello" -}}
      {{ $greeting }}, {{ @name }}!
      {{ if @age -}}
      You are {{ @age }} years old.
      {{- else -}}
      Age unknown.
      {{- end }}
      Your friends:
      {{- range @friends }}
      - {{ . }}
      {{- else }}
      You have no friends.
      {{- end }}
      Profile details:
      {{- range $key, $value := @profile }}
      - {{ $key }}: {{ $value }}
      {{- end }}
      {{ define "footer" -}}
      This is the footer for {{ @name }}.
      {{- end }}
      {{ template "footer" . -}}
      ------------------------------------
      Content-Type: {{ #Content_Type }}
      Redacted Auth Header: {{ #Authorization }}
```

> [!NOTE]
> Supported [placeholder types](../usage/advanced#placeholders):
>
> | `.` Variables | `@` Body | `#` Headers |
> | ------------- | -------- | ----------- |
> | ✅            | ✅       | ✅          |
