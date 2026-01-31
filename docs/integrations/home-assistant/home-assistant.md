---
title: Home Assistant
---

# Home Assistant

Instructions on how you can use **Secured Signal API** as a notification service for [Home Assistant](https://github.com/home-assistant/core).

## Setup

### 1. Home Assistant Configuration

To be able to use the Signal Messenger integration in Home Assistant you need to modify or add the following to your `configuration.yml` file:

```yaml
{{{ #://./configuration.yml }}}
```

> [!TIP]
> If you want to use a recipients placeholder (_array_) instead of _single_ recipients, modify the `url` to `http://sec-signal-api/@auth=API_TOKEN/@recipients={{.RECIPIENTS}}`

Here we are taking advantage of the `url` field for adding `/@auth=API_TOKEN` in order to use [Path Auth](../usage#auth).

For more detailed configuration instructions read the [official Home Assistant docs](https://www.home-assistant.io/integrations/signal_messenger/).

### 2. Enabling Path Auth

By default, [Path Auth](../usage#auth) is disabled, so we first need to enable it in the config by adding `path` to [`auth.methods`](../configuration/auth):

```yaml
api:
  auth:
    methods: [bearer, basic, body, path]
```

And that's basically it, have fun!

## Sources

- https://www.home-assistant.io/integrations/signal_messenger/
