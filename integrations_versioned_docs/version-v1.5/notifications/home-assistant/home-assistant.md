---
title: Home Assistant
sidebar_custom_props:
  icon: https://upload.wikimedia.org/wikipedia/commons/a/ab/New_Home_Assistant_logo.svg
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
> If you want to use a list of recipients as placeholder you have to add `/@recipients={{.RECIPIENTS}}` to the `url`:
>
> ```yaml
> url: "http://api:API_TOKEN@sec-signal-api:8880/@recipients={{.RECIPIENTS}}"
> ```

Here we are taking advantage of the `url` field to implement [Basic Auth](/docs/usage#auth) by using `user:password@host:port`.

For more detailed configuration instructions read the [official Home Assistant docs](https://www.home-assistant.io/integrations/signal_messenger/).

And that's basically it, have fun!

## Sources

- https://www.home-assistant.io/integrations/signal_messenger/
