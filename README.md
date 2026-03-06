<img align="center" width="1048" height="512" alt="Secure API Gateway Proxy for Signal CLI REST API" src="https://github.com/codeshelldev/secured-signal-api/raw/refs/heads/docs/static/img/banner.png" />

<h3 align="center">Secure API Gateway Proxy for <a href="https://github.com/bbernhard/signal-cli-rest-api">Signal CLI REST API</a></h3>

<p align="center">
token-based authentication,
endpoint restrictions, placeholders, flexible configuration
</p>

<p align="center">
🔒 Secure · ⭐️ Configurable · 🚀 Easy to Deploy with Docker
</p>

<div align="center">
  <a href="https://github.com/codeshelldev/secured-signal-api/releases">
    <img 
		src="https://img.shields.io/github/v/release/codeshelldev/secured-signal-api?sort=semver&logo=github&label=Release" 
		alt="GitHub release"
	>
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/stargazers">
    <img 
		src="https://img.shields.io/github/stars/codeshelldev/secured-signal-api?style=flat&logo=github&label=Stars" 
		alt="GitHub stars"
	>
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/pkgs/container/secured-signal-api">
    <img 
		src="https://ghcr-badge.egpl.dev/codeshelldev/secured-signal-api/size?color=%2344cc11&tag=latest&label=Image+Size&trim="
		alt="Docker image size"
	>
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/pkgs/container/secured-signal-api">
    <img 
		src="https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fghcr-badge.elias.eu.org%2Fapi%2Fcodeshelldev%2Fsecured-signal-api%2Fsecured-signal-api&query=downloadCount&label=Pulls&color=2344cc11&logo=docker"
		alt="Docker image Pulls"
	>
  </a>
  <a href="./LICENSE">
    <img 
		src="https://img.shields.io/badge/License-MIT-green.svg"
		alt="License: MIT"
	>
  </a>
</div>

## Contents

> [!IMPORTANT]
> Check out the [**Official Documentation**](https://codeshelldev.github.io/secured-signal-api) for up-to-date instructions and additional content!

> [!WARNING]
> We are slowly moving away from this README and instead are trying to make the [**Official Documentation**](https://codeshelldev.github.io/secured-signal-api) the only source of truth

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Features](#features)
- [Configuration](#configuration)
  - [Endpoints Restrictions](https://codeshelldev.github.io/secured-signal-api/docs/configuration/endpoints)
  - [Placeholder Variables](https://codeshelldev.github.io/secured-signal-api/docs/configuration/variables)
  - [Field Policies](https://codeshelldev.github.io/secured-signal-api/docs/configuration/field-policies)
  - [Field Mappings](https://codeshelldev.github.io/secured-signal-api/docs/configuration/field-mappings)
  - [Message Template](https://codeshelldev.github.io/secured-signal-api/docs/configuration/message-template)
  - [Port Restrictions](https://codeshelldev.github.io/secured-signal-api/docs/configuration/port)
  - [Hostname Restrictions](https://codeshelldev.github.io/secured-signal-api/docs/configuration/hostnames)
  - [IP Filter](https://codeshelldev.github.io/secured-signal-api/docs/configuration/ip-filter)
  - [Rate Limiting](https://codeshelldev.github.io/secured-signal-api/docs/configuration/rate-limiting)
  - [Trusted Proxies](https://codeshelldev.github.io/secured-signal-api/docs/configuration/trusted-proxies)
  - [Trusted IPs](https://codeshelldev.github.io/secured-signal-api/docs/configuration/trusted-ips)
  - [Log Level](https://codeshelldev.github.io/secured-signal-api/docs/configuration/log-level)
  - [Port](https://codeshelldev.github.io/secured-signal-api/docs/configuration/port)
  - [Hostnames](https://codeshelldev.github.io/secured-signal-api/docs/configuration/hostnames)
  - [Auth Methods](https://codeshelldev.github.io/secured-signal-api/docs/configuration/auth)
- [Reverse Proxy](https://codeshelldev.github.io/secured-signal-api/docs/reverse-proxy)
- [Integrations](https://codeshelldev.github.io/secured-signal-api/docs/integrations)
- [Contributing](#contributing)
- [Support](#support)
- [Help](#help)
- [License](#license)
- [Legal](#legal)

## Getting Started

### Installation

---
sidebar_position: 2
title: Installation
---

# Installation

Get the latest version of the `docker-compose.yaml` file:

```yaml
file not found: /home/runner/work/secured-signal-api/secured-signal-api/.github/templates/docs/getting-started/examples/docker-compose.yaml
```

> [!IMPORTANT]
> In this documentation, we'll be using `sec-signal-api:8880` as the host for simplicity,
> please replace it with your actual container/host IP, port, or hostname

## API Tokens

Now head to [configuration](../configuration/api-tokens) and define some **API tokens**.

> [!TIP]
> This recommendation is part of the [**best practices**](../best-practices)

### Setup

Once you have installed **Secured Signal API** you can [register or link a signal account](https://codeshelldev.github.io/secured-signal-api/docs/getting-started/setup).

## Usage

---
sidebar_position: 1
title: Usage
---

# Usage

In this section we'll be taking a look at how to use **Secured Signal API**.

## Basic

Here is a quick command to see if you've correctly followed the [setup instructions](./getting-started/setup):

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"message":"Hello, World!", "number":"<from>", "recipients":["<to>"]}' \
    'http://sec-signal-api:8880/v2/send'
```

This will send `Hello, World!` to `<to>` from `<from>`.

## Auth

**Secured Signal API** implements 5 auth methods:

| Method      | Example                                                                                                      |
| :---------- | :----------------------------------------------------------------------------------------------------------- |
| Bearer Auth | `Authorization: Bearer API_TOKEN` (header)                                                                   |
| Basic Auth  | `Authorization: Basic base64(api:API_TOKEN)` (header)<br/>`http://api:API_TOKEN@host:port` (client specific) |
| Query Auth  | `http://host:port/abc?@auth=API_TOKEN` (query parameter)                                                     |
| Path Auth   | `http://host:port/@auth=API_TOKEN/abc` (path parameter)                                                      |
| Body Auth   | `{ "auth": "API_TOKEN" }` (request body field)                                                               |

> [!WARNING]
> **Query** and **Path** auth are disabled by default and [must be enabled in the config](../configuration/auth.md)

**Example:**

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer API_TOKEN" \
    -d '{"message":"Hello, World!", "number":"<from>", "recipients":["<to>"]}' \
    'http://sec-signal-api:8880/v2/send'
```

## Features

---
sidebar_position: 2
title: Features
---

# Features

Here are some of the highlights of using **Secured Signal API**.

## Message Template

> _Structure your messages_

**Message Templates** can be used to customize your final message after preprocessing.
Look at this complex template for example:

```yaml
file not found: /home/runner/work/secured-signal-api/secured-signal-api/.github/templates/docs/configuration/examples/message-template.yml
```

It can extract needed data from the body and headers to then process them using Go's templating library
and finally output a message packed with so much information.

Head to [Configuration](./configuration/templating#message-template) to see how-to use.

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

## Rate Limiting

**Rate Limiting** is used for limiting requests and to stop server overload, because of DDoS attacks, malconfigured clients, or malicious actors.  
It ensures fair usage per token by controlling how many requests can be processed within a defined period.

Limit those rates [here](./configuration/rate-limiting).

## Endpoints

> _Block unwanted access_

**Endpoints** are used for restricting unauthorized access and for ensuring least privilege.

[Let's start blocking then!](./configuration/endpoints)

## IP Filters

**IP Filters** are used for restricting access to **Secured Signal API** by blocking or specifically allowing IPs and CIDR ranges.

Configure your _mini firewall_ [here](./configuration/ip-filter).

## Configuration

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

```md
config.yml
tokens
├── notify.yml
└── totp.yml
```

Config files are the **recommended** way to configure and use **Secured Signal API**,
they are **flexible**, **extensible** and really **easy to use**.

> [!TIP]
> Configs also support placeholders, for example:
> `${{ .env.NUMBER }}` or `${{ .vars.RECIPIENTS }}`
>
> - Use `.vars` for placeholders from [variables](./variables)
> - and `.env` for environment variables

> [!NOTE]
> To change the internal config file location set `CONFIG_PATH` in your **environment** to an absolute path (default: `/config/config.yml`)

This example config shows all the individual settings that can be applied:

```yaml
file not found: /home/runner/work/secured-signal-api/secured-signal-api/.github/templates/docs/configuration/examples/config.yml
```

### Token Configs

```
tokens
├── notify.yml
└── totp.yml
```

> But wait! There is more… 😁

Token configs can be used to create **per-token** defined **overrides** and settings.

> [!NOTE]
> Create them under `TOKENS_PATH` (default: `config/tokens/`)

Here is an example:

```yaml
file not found: /home/runner/work/secured-signal-api/secured-signal-api/.github/templates/docs/configuration/examples/token.yml
```

## Contributing

Found a bug? Want to change or add something?
Feel free to open up an [issue](https://github.com/codeshelldev/secured-signal-api/issues) or create a [pull request](https://github.com/codeshelldev/secured-signal-api/pulls)!

## Support

Has this Repo been helpful 👍️ to you? Then consider ⭐️'ing this Project.

:)

## Help

**Are you having problems setting up Secured Signal API?**<br/>
No worries check out the [discussions](https://github.com/codeshelldev/secured-signal-api/discussions) tab and ask for help.

**We are all volunteers**, so please be friendly and patient.

## License

This Project is licensed under the [MIT License](./LICENSE).

## Legal

Logo designed by [@CodeShellDev](https://github.com/codeshelldev), All Rights Reserved.

This Project is not affiliated with the Signal Foundation.
