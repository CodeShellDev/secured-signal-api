
<img align="center" width="1048" height="512" alt="Secure API Gateway Proxy for Signal CLI REST API" src="https://github.com/codeshelldev/secured-signal-api/raw/refs/heads/docs/static/img/brand/banner.png" />

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
		alt="GitHub release">
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/stargazers">
    <img 
		src="https://img.shields.io/github/stars/codeshelldev/secured-signal-api?style=flat&logo=github&label=Stars" 
		alt="GitHub stars">
  </a>
  <a href="https://matrix.to/#/#secured-signal-api:matrix.org">
	  <img 
		  alt="Matrix" 
		  src="https://img.shields.io/matrix/secured-signal-api%3Amatrix.org?style=flat&color=blue&logo=matrix&label=Chat">
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/pkgs/container/secured-signal-api">
    <img
    src='https://img.shields.io/badge/Image%20Size-<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="136" height="20">
    <linearGradient id="b" x2="0" y2="100%">
        <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
        <stop offset="1" stop-opacity=".1"/>
    </linearGradient>
    <mask id="anybadge_36190">
        <rect width="136" height="20" rx="3" fill="#fff"/>
    </mask>
    <g mask="url(#anybadge_36190)">
        <path fill="#555" d="M0 0h75v20H0z"/>
        <path fill="#44cc11" d="M75 0h61v20H75z"/>
        <path fill="url(#b)" d="M0 0h136v20H0z"/>
    </g>
    <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
        <text x="38.5" y="15" fill="#010101" fill-opacity=".3">Image Size</text>
        <text x="37.5" y="14">Image Size</text>
    </g>
    <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
        <text x="106.5" y="15" fill="#010101" fill-opacity=".3">9.81 MiB</text>
        <text x="105.5" y="14">9.81 MiB</text>
    </g>
</svg>
-_?color=2344cc11'
    alt="Docker image Pulls">
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/pkgs/container/secured-signal-api">
    <img 
		src="https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fghcr-badge.elias.eu.org%2Fapi%2Fcodeshelldev%2Fsecured-signal-api%2Fsecured-signal-api&query=downloadCount&label=Pulls&color=2344cc11&logo=docker"
		alt="Docker image Pulls">
  </a>
  <a href="./LICENSE">
    <img 
		src="https://img.shields.io/badge/License-MIT-green.svg"
		alt="License: MIT">
  </a>
</div>

## Contents

> [!IMPORTANT]
> Check out the [**Official Documentation**](https://codeshelldev.github.io/secured-signal-api) for up-to-date instructions and additional content!

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
  - [CORS](https://codeshelldev.github.io/secured-signal-api/docs/configuration/cors)
- [Reverse Proxy](https://codeshelldev.github.io/secured-signal-api/docs/reverse-proxy)
- [Integrations](https://codeshelldev.github.io/secured-signal-api/docs/integrations)
- [Contributing](#contributing)
- [Support](#support)
- [Help](#help)
- [License](#license)
- [Legal](#legal)

## Getting Started

### Installation

Get the latest version of the `docker-compose.yaml` file:

```yaml
services:
  signal-api:
    image: bbernhard/signal-cli-rest-api:latest
    container_name: signal-api
    environment:
      - MODE=normal
    volumes:
      - ./data:/home/.local/share/signal-cli
    restart: unless-stopped
    networks:
      backend:
        aliases:
          - signal-api

  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api:latest
    container_name: secured-signal
    environment:
      API__URL: http://signal-api:8080
      SETTINGS__MESSAGE__VARIABLES__RECIPIENTS: "[+123400002, +123400003, +123400004]"
      SETTINGS__MESSAGE__VARIABLES__NUMBER: "+123400001"
      API__TOKENS: "[LOOOOOONG_STRING]"
    ports:
      - "8880:8880"
    volumes:
      - ./db:/db
    restart: unless-stopped
    networks:
      backend:
        aliases:
          - secured-signal-api


networks:
  backend:
```

> [!IMPORTANT]
> In this documentation, we'll be using `sec-signal-api:8880` as the host for simplicity,
> please replace it with your actual container/host IP, port, or hostname

#### API Tokens

Now head to [configuration](https://codeshelldev.github.io/secured-signal-api/docs/configuration/api-tokens) and define some **API tokens**.

> [!TIP]
> This recommendation is part of the [**best practices**](https://codeshelldev.github.io/secured-signal-api/docs/best-practices)

### Setup

Once you have installed **Secured Signal API** you can [register or link a signal account](https://codeshelldev.github.io/secured-signal-api/docs/getting-started/setup).

## Usage

### Basic

Here is a quick command to see if you've correctly followed the [setup instructions](https://codeshelldev.github.io/secured-signal-api/docs/getting-started/setup):

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"message":"Hello, World!", "number":"<from>", "recipients":["<to>"]}' \
    'http://sec-signal-api:8880/v2/send'
```

This will send `Hello, World!` to `<to>` from `<from>`.

### Auth

**Secured Signal API** implements 5 auth methods:

| Method      | Example                                                                                                      |
| :---------- | :----------------------------------------------------------------------------------------------------------- |
| Bearer Auth | `Authorization: Bearer API_TOKEN` (header)                                                                   |
| Basic Auth  | `Authorization: Basic base64(api:API_TOKEN)` (header)<br/>`http://api:API_TOKEN@host:port` (client specific) |
| Query Auth  | `http://host:port/abc?@auth=API_TOKEN` (query parameter)                                                     |
| Path Auth   | `http://host:port/@auth=API_TOKEN/abc` (path parameter)                                                      |
| Body Auth   | `{ "auth": "API_TOKEN" }` (request body field)                                                               |

> [!WARNING]
> **Query** and **Path** auth are disabled by default and [must be enabled in the config](https://codeshelldev.github.io/secured-signal-api/docs/configuration/auth)

**Example:**

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer API_TOKEN" \
    -d '{"message":"Hello, World!", "number":"<from>", "recipients":["<to>"]}' \
    'http://sec-signal-api:8880/v2/send'
```

## Features

### Message Template

> _Structure your messages_

**Message Templates** can be used to customize your final message after preprocessing.
Look at this complex template for example:

```yaml
settings:
  message:
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

```

It can extract needed data from the body and headers to then process them using Go's templating library
and finally output a message packed with so much information.

Head to [Configuration](https://codeshelldev.github.io/secured-signal-api/docs/configuration/message-template) to see how-to use.

### Placeholders

> _Time saving and flexible_

**Placeholders** are one of the highlights of Secured Signal API,
these have saved me, and will save many others, much time by, for example, not having to change your phone number in every service separately.

Take a look at the [usage](https://codeshelldev.github.io/secured-signal-api/docs/usage/advanced).

### Field Mappings

> _Standardize output_

**Field Mappings** are very useful for when your favorite service does not officially support **Secured Signal API** (or Signal CLI REST API).
With this feature you have the power to do it yourself, just extract what's needed and then integrate with any of the other features.

Interested? [Take a look](https://codeshelldev.github.io/secured-signal-api/docs/configuration/field-mappings).

### Field Policies

**Field Policies** are a great way to disallow specific fields or even disallowing fields with unwanted values.
This is really helpful when trying to block certain numbers from using certain tokens, and therefor mitigating risks of unwanted use of a token.

Find more about this feature [here](https://codeshelldev.github.io/secured-signal-api/docs/configuration/field-policies).

### Rate Limiting

**Rate Limiting** is used for limiting requests and to stop server overload, because of DDoS attacks, malconfigured clients, or malicious actors.  
It ensures fair usage per token by controlling how many requests can be processed within a defined period.

Limit those rates [here](https://codeshelldev.github.io/secured-signal-api/docs/configuration/rate-limiting).

### Endpoints

> _Block unwanted access_

**Endpoints** are used for restricting unauthorized access and for ensuring least privilege.

[Let's start blocking then!](https://codeshelldev.github.io/secured-signal-api/docs/configuration/endpoints)

### IP Filters

**IP Filters** are used for restricting access to **Secured Signal API** by blocking or specifically allowing IPs and CIDR ranges.

Configure your _mini firewall_ [here](https://codeshelldev.github.io/secured-signal-api/docs/configuration/ip-filter).

### CORS

> _Enable secure browser access_

**CORS** support allows web applications running in a browser to access the API directly from a different origin.

Learn more about browser integration [here](https://codeshelldev.github.io/secured-signal-api/docs/configuration/cors).

## Configuration

### Environment Variables

Whilst being a bit **restrictive** environment variables are a great way to configure Secured Signal API.

Suppose you want to set a new [placeholder](https://codeshelldev.github.io/secured-signal-api/docs/usage/advanced#placeholders) `NUMBER` in your environment…

```yaml
environment:
  SETTINGS__MESSAGE__VARIABLES__NUMBER: "+123400001"
```

This would internally be converted into `settings.message.variables.number` matching the config formatting.

> [!IMPORTANT]
> Single underscores `_` are removed during conversion, whereas double underscores `__` convert the variable into a nested object (with `__` replaced by `.`)

### Config Files

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
> - Use `.vars` for placeholders from [variables](https://codeshelldev.github.io/secured-signal-api/docs/variables)
> - and `.env` for environment variables

> [!NOTE]
> To change the internal config file location set `CONFIG_PATH` in your **environment** to an absolute path (default: `/config/config.yml`)

This example config shows all the individual settings that can be applied:

```yaml
## Example Config (all configurations shown)
service:
  logLevel: info
  port: 8880
  hostnames:
    - mydomain.com

api:
  url: http://signal-api:8080
  tokens: [token1, token2]
  auth:
    methods: [bearer, basic, body]
    tokens:
      - set: [pathToken1, pathToken2]
        methods: [path]
      - set: [queryAndBodyToken]
        methods: [body, query]

settings:
  http:
    responseHeaders:
      X-Custom: "xyz"

  message:
    messageTemplate: |
      You've got a Notification:
      {{ @message }} 
      At {{ @data.timestamp }} on {{ @data.date }}.
      Send using {{ .NUMBER }}.

    templating:
      body: true
      query: true
      path: true

    scheduling:
      enabled: true
      maxHorizon: 10d

    injecting:
      urlToBody:
        query: true
        path: true

    variables:
      number: "+123400001"
      recipients: ["+123400002", "group.id", "user.id"]
      text_mode: "${{ .env.TEXT_MODE }}"

    fieldMappings:
      "@message":
        - field: "msg"
          score: 100

  access:
    trustedIPs:
      - 192.168.1.10

    trustedProxies:
      - 172.20.0.100

    ipFilter:
      allowed:
        - 192.168.1.10
        - 192.168.2.0/24
      blocked:
        - 192.168.2.44

    endpoints:
      allowed:
        - pattern: /v2/send
        - pattern: /v1/receive/${{ .vars.NUMBER }}
      blocked:
        - pattern: /v1/about
        - pattern: /v1/receive/[^0]+
          matchType: regex

    rateLimiting:
      limit: 100
      period: 1h

    fieldPolicies:
      "@number":
        - value: "+123400003"
          action: block
        - value: "+12340000[4-9]"
          matchType: regex
          action: block

    cors:
      methods: [GET, POST]
      headers: ["Content-Type", "Accept-Language"]
      origins:
        - url: "https://domain.com"
          methods: [GET]
          headers: ["Content-Type"]

```

#### Token Configs

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
## Example Token Config (overwrites)
service:
  logLevel: info
  port: 8880
  hostnames:
    - mydomain.com

api:
  tokens: [token1, token2]
  auth:
    methods: [bearer, basic, body, path] # add path auth

settings:
  http:
    responseHeaders: # overwrite response headers
      X-Custom: "Lorem Ipsum Dolor"

  message:
    scheduling: # disable
    injecting: # disable
    templating: # disable
    messageTemplate: # disable

    variables: # overwrite main config variables
      number: "+123400010"
      recipients: ["+123400020", "group.id", "user.id"]

    fieldMappings: # overwrite @message from main config
      "@message":
        - field: "msg"
          score: 100

  access:
    trustedIPs: # disable
    trustedProxies: # disable
    ipFilter: # disable
    fieldPolicies: # disable
    cors: # disable

    endpoints: # overwrite main config endpoints
      allowed:
        - pattern: /v1/receive/${{ .vars.NUMBER }}
        - pattern: /v2/send
      blocked:
        - pattern: /v1/about

    rateLimiting:
      limit: 100
      period: 10h # overwrite main config period

```

## Contributing

Found a bug? Want to change or add something?
Feel free to open up an [issue](https://github.com/codeshelldev/secured-signal-api/issues) or create a [pull request](https://github.com/codeshelldev/secured-signal-api/pulls)!

## Support

Has this Repo been helpful 👍️ to you? Then consider ⭐️'ing this Project.

:)

## Help

**Are you having problems setting up Secured Signal API?**<br/>
No worries check out the [discussions](https://github.com/codeshelldev/secured-signal-api/discussions) tab or our [matrix chat](https://matrix.to/#/#secured-signal-api:matrix.org) and ask for help.

**We are all volunteers**, so please be friendly and patient.

## License

This Project is licensed under the [MIT License](./LICENSE).

## Legal

Logo designed by [@CodeShellDev](https://github.com/codeshelldev), All Rights Reserved.

This Project is not affiliated with the Signal Foundation.
