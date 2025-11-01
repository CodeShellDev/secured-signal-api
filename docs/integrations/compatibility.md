---
sidebar_position: 1
title: Compatibility
---

# Compatibility

## The Problem

**Secured Signal API** is only of use when compatible with programms.
Even if it keeps the underlying Signal CLI REST API you'd still need your services to support Signal CLI REST API.

## The Solution

**Secured Signal API** implements enough features to technically support any and all services.
But with one flaw:

> _manual configuration_

In order for Secured Signal API to be compatible and integratable with a service you still need to manually define [**Field Mappings**](../configuration/field-mappings)
and [**Message Templates**](../configuration/message-template), which is quite easy,
provided you know what the services is using as payload (try sending a request to some debugging endpoint).

> _Now wouldn't it be great if someone had already done that?_

If you are using a common and popular service or programm there is probably someone who already configured everything and was willing to share it on
[our Github Discussions](https://github.com/codeshelldev/secured-signal-api/discussions) (**Thank you!**).

## How to Help

You successfully integrated a service and want to share it?

> Well that's nice of you 🤩👍️

Then create a [Discussion](https://github.com/CodeShellDev/secured-signal-api/discussions/categories/integrations) and share your configs or if you want you can also submit a [Pull Request](https://github.com/codeshelldev/secured-signal-api/pulls) to add your integration to the **Integrations Section** in the official Documentation.
