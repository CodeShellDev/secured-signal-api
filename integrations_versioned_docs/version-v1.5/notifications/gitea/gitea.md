---
title: Gitea
sidebar_custom_props:
  icon: https://raw.githubusercontent.com/go-gitea/gitea/refs/heads/main/assets/logo.svg
---

# Gitea

Here's how you can use **Secured Signal API** as a notification service for [Gitea](https://github.com/go-gitea/gitea).

## Setup

### 1. Message Template

Because Gitea's webhook data is very _clustered_, we need to use [**Message Templates**](../configuration/message-template) to ensure correct message rendering.

Here is an example:

```yaml
{{{ #://./message-template.yml }}}
```

Add this to your token config and modify it to your needs.

### 2. Webhook

Head to your Gitea repository (or user settings) and go to `Settings > Webhooks` and create a new Gitea webhook.

![Webhook](/integrations/gitea/webhook.png)

## Testing

After you've completed the setup you can try out your new notification integration:

![Example Issue](/integrations/gitea/issue.png)

```markdown
📝 **#1 Very Important Issue**  
🟢 | 👤 User
🔗 https://gitea.domain.com/user/repo/issues/1
```

## Features

The provided Message Template currently supports:

- Gitea issues
- Git commits
- Gitea pull requests
