---
sidebar_position: 2
title: Features
---

# Features

Here are some of the highlights of using **Secured Signal API**

## Message Template

> _Structure your messages_

**Message Templates** can be used to customize your final message after preprocessing.
Look at this complex template for example:

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

It can extract needed data from the body and headers to then process them using Go's templating library
and finally output a message packed with so much information.

Head to [configuration](./configuration/message-template) to see how-to use.

## Placeholders

> _Time saving and flexible_

**Placeholders** are one of the highlights of Secured Signal API,
these have saved me, and will save many others, much time by, for example, not having to change your phone number in every service separately.

Take a look at the [usage](./usage/advanced).

## Field Mappings

> _Standardize output_

**Field Mappings** are very useful for when your favorite service does not officially support **Secured Signal API** (or Signal CLI REST API).
With this feature you have the power to do it yourself, just extract what's needed and then integrate with any of the other features.

Interested? [Take a look](./configuration/field-mappings).

## Field Policies

**Field Policies** are a great way to disallow specific fields or even disallowing fields with unwanted values.
This is really helpful when trying to block certain numbers from using certain tokens, and therefor mitigating risks of unwanted use of a token.

Find more about this feature [here](./configuration/field-policies).

## Endpoints

> _Block unwanted access_

**Endpoints** are used for restricting unauthorized access and for ensuring least privilege.
[Let's start blocking then!](./configuration/endpoints)
