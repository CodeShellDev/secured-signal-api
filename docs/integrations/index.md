---
sidebar_position: 1
title: Compatibility
---

# Compatibility

## The Problem

**Secured Signal API** is only of use when compatible with programs.
Even tho it keeps the underlying Signal CLI REST API you'd still need your services to support Signal CLI REST API.

## The Solution

**Secured Signal API** implements enough features to technically support any and all services.
But with one flaw:

> _manual configuration_

In order for **Secured Signal API** to be compatible and integratable with a service, you still need to manually define [**Field Mappings**](./configuration/field-mappings) and [**Message Templates**](./configuration/message-template).

This process is straightforward, provided you know what the service uses as its payload — for example, you can test by sending a request to a debugging endpoint.

> _Now wouldn't it be great if someone had already done that?_

If you are using a common and popular service or programs there is probably someone who already configured everything and was willing to share it on
[our GitHub discussions](https://github.com/codeshelldev/secured-signal-api/discussions) (**Thank you!**).

## How to Help

You successfully integrated a service and want to share it?

> Well, that's nice of you 🤩👍️

Then create a [discussion](https://github.com/CodeShellDev/secured-signal-api/discussions/categories/integrations) and share your configs or if you want you can submit a [pull request](https://github.com/codeshelldev/secured-signal-api/pulls) to add your integration to the **integrations section** in the official documentation.
