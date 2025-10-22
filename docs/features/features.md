---
sidebar_position: 7
title: Features
---

# Features

Here are some of the highlights of using **Secured Signal API**

---

## Message Templates

> _Incredible fun and useful_

**Message Templates** are used to customize your final message after preprocessing.
Look at this complex template for example:

```yaml
settings:
    messageTemplate: |
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

It can extract needed data from the body and even the headers (with an exception of the `Authorization` Header) and then process them using Go's Templating Library
and finally output a message packed with so much information.

Head to [Configuration](./configuration/message-templates) to see how-to use.

---

## Placeholders

> _Timesaving and flexible_

**Placeholders** are one of the highlights of this Project,
these have saved me and will save many others much time by not having to change your phone number in every service separately or other values.

Take a look at the [Usage](./usage/advanced).

---

## Data Aliases

> _Boring, but sooo definetly needed_

**Data Aliases** are also very useful for when your favorite service does not officially support **Secured Signal API** (or Signal CLI REST API).
With this feature you have the power to do it yourself, just extract what's needed and then integrate with any of the other features.

Interested? [Take a look](./configuration/data-aliases).

---

## Endpoints

> _why do you need write access for reading messages?!_

**Endpoints** or rather their subfeatures:

- [**Allowed Endpoints**](./configuration/endpoints)
- [**Blocked Endpoints**](./configuration/endpoints)

Go hand in hand for restricting unauthorized access and for ensuring least privilege.
[Time to go blocking...](./configuration/endpoints)

---
