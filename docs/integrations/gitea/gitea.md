---
title: Gitea
---

# Gitea

Here's how you can use **Secured Signal API** as a notification service for [gitea](https://github.com/go-gitea/gitea).

## Setup

### 1. Message Template

Because gitea's webhook data is very _clustered_, we need to use [**Message Templates**](../configuration/message-template) to ensure correct message rendering.

Here is an example:

```yaml
settings:
  message:
    template: |
      {{- if and @issue (ne @issue nil) (not @is_pull) -}}
      📝 **#{{@issue.number}} {{@issue.title}}**  
      {{ if eq @issue.state "open" -}}🟢{{- else if eq @issue.state "closed" -}}🔴{{- else -}}{{@issue.state}}{{- end }} | 👤 {{@sender.full_name}}
      🔗 {{@issue.html_url}}
      {{- end -}}

      {{- if and @pull_request (ne @pull_request nil) -}}
      🚀 **#{{@pull_request.number}} {{@pull_request.title}}**  
      {{ if eq @pull_request.state "open" -}}🟢{{- else if eq @pull_request.state "closed" -}}🔴{{- else if eq @pull_request.state "merged" -}}🟣{{- else -}}{{@pull_request.state}}{{- end }} | 👤 {{@sender.full_name}}
      🔗 {{@pull_request.html_url}}
      {{- end -}}

      {{- if and @commits (gt (len @commits) 0) }}
      📥️ **Push** → *{{@ref}}*  
      📁 {{@repository.full_name}} | 👤 {{@pusher.full_name}} | 🔢 {{@total_commits}}  
      {{- range @commits }}
      - 🧾 *{{.message}}*  
        {{- if .added }}
          {{- range .added }}➕ {{.}} {{- end }}
        {{- end -}}
        {{- if .modified }}
          {{- range .modified }} ✏️ {{.}} {{- end }}
        {{- end -}}
        {{- if .removed }}
          {{- range .removed }} ❌ {{.}} {{- end }}
        {{- end }}
        🔗 {{.url}}
      {{- end }}

      🔎 Compare: {{@compare_url}}
      {{- end -}}
```

Add this to your token config and modify it to your needs.

### 2. Webhook

Head to your gitea repository (or user settings) and go to `Settings > Webhooks` and create a new Gitea webhook.

![Webhook](/integrations/gitea/webhook.png)

## Testing

After you've completed the Setup you can try out your new notification integration:

![Example Issue](/integrations/gitea/issue.png)

```markdown
📝 **#1 Very Important Issue**  
🟢 | 👤 User
🔗 https://gitea.domain.com/user/repo/issues/1
```

## Features

The provided Message Template currently supports:

- issues
- commits
- pull requests
