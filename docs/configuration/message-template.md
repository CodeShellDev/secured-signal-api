---
title: Message Template
---

# Message Template

Message Templates are the best way to **structure** and **customize** your messages and also very useful for **compatiblity** between different services.

Configure them by using the `message.template` attribute in you config.

These support Go Templates (See [Templates](../usage/formatting)) and work by templating the `message` attribute in the request's body.

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

> [!IMPORTANT]
> Message Templates support [Standard Golang Templating](../usage/formatting).
> Use `@data.key` to reference Body Keys, `#Content_Type` for Headers and `.KEY` for [Variables](./variables).

> [!WARNING]
> Templating using the `Authorization` header results in a redacted string.
