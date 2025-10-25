---
sidebar_position: 1
title: Overview
---

# Configuration

Here is how you configure **Secured Signal API**

## Environment Variables

While being a bit **restrictive** environment variables are a great way to configure Secured Signal API.

Suppose you want to set a new [Placeholder](../usage/advanced) `NUMBER` in your Environment...

```yaml
environment:
  SETTINGS__MESSAGE__VARIABLES__NUMBER: "+123400001"
```

This would internally be converted into `settings.message.variables.number` matching the config formatting.

> [!IMPORTANT]
> Underscores `_` are removed during Conversion, double Underscores `__` on the other hand convert the Variable into a nested Object (`__` replaced by `.`)

## Config Files

Config files are the **recommended** way to configure and use **Secured Signal API**,
they are **flexible**, **extensible** and really **easy to use**.

Config files allow **YAML** formatting and also `${ENV}` to get environment variables.

> [!NOTE]
> To change the internal config file location set `CONFIG_PATH` in your **Environment** to an absolute path. (default: `/config/config.yml`)

This example config shows all of the individual settings that can be applied:

```yaml
{{{ #://./examples/config.yml }}}
```

### Token Configs

> But wait! There is more... 😁

Token Configs are used to create **per-toke**n defined **overrides** and settings.

> [!NOTE]
> Create them under `TOKENS_PATH` (default: `config/tokens/`)

This way you can permission tokens by further restricting or adding [Endpoints](../configuration/endpoints), [Placeholders](../configuration/variables), etc.

Here is an example:

```yaml
{{{ #://./examples/token.yml }}}
```
