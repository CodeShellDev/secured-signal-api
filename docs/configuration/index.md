---
sidebar_position: 1
title: Configuration
---

# Configuration

Here is how you configure **Secured Signal API**

## Environment Variables

Whilst being a bit **restrictive** environment variables are a great way to configure Secured Signal API.

Suppose you want to set a new [placeholder](./usage/advanced#placeholders) `NUMBER` in your environment…

```yaml
environment:
  SETTINGS__MESSAGE__VARIABLES__NUMBER: "+123400001"
```

This would internally be converted into `settings.message.variables.number` matching the config formatting.

> [!IMPORTANT]
> Single underscores `_` are removed during conversion, whereas double underscores `__` convert the variable into a nested object (with `__` replaced by `.`)

## Config Files

Config files are the **recommended** way to configure and use **Secured Signal API**,
they are **flexible**, **extensible** and really **easy to use**.

Config files allow **YAML** formatting and additionally `${ENV}` to get environment variables.

> [!TIP]
> To change the internal config file location set `CONFIG_PATH` in your **environment** to an absolute path (default: `/config/config.yml`)

This example config shows all the individual settings that can be applied:

```yaml
{{{ #://./examples/config.yml }}}
```

### Token Configs

> But wait! There is more… 😁

Token configs can be used to create **per-token** defined **overrides** and settings.

> [!IMPORTANT]
> Create them under `TOKENS_PATH` (default: `config/tokens/`)

Here is an example:

```yaml
{{{ #://./examples/token.yml }}}
```
