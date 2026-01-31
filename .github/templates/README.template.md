<img align="center" width="1048" height="512" alt="Secure Proxy for Signal CLI REST API" src="https://github.com/codeshelldev/secured-signal-api/raw/refs/heads/docs/static/img/banner.png" />

<h3 align="center">Secure Proxy for <a href="https://github.com/bbernhard/signal-cli-rest-api">Signal CLI REST API</a></h3>

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
- [Setup](#setup)
- [Usage](#usage)
- [Features](https://codeshelldev.github.io/secured-signal-api/docs/features)
- [Configuration](#configuration)
  - [Endpoints](#endpoints)
  - [Variables](#variables)
  - [Field Policies](#field-policies)
  - [Field Mappings](#field-mappings)
  - [Message Templates](#message-templates)
  - [Port Restrictions](https://codeshelldev.github.io/secured-signal-api/docs/configuration/port)
  - [Hostname Restrictions](https://codeshelldev.github.io/secured-signal-api/docs/configuration/hostnames)
  - [IP Filter](https://codeshelldev.github.io/secured-signal-api/docs/configuration/ip-filter)
  - [Rate Limiting](https://codeshelldev.github.io/secured-signal-api/docs/configuration/rate-limiting)
  - [Trusted Proxies](https://codeshelldev.github.io/secured-signal-api/docs/configuration/trusted-proxies)
  - [Trusted IPs](https://codeshelldev.github.io/secured-signal-api/docs/configuration/trusted-ips)
  - [Log Level](https://codeshelldev.github.io/secured-signal-api/docs/configuration/log-level)
  - [Port](https://codeshelldev.github.io/secured-signal-api/docs/configuration/port)
  - [Hostnames](https://codeshelldev.github.io/secured-signal-api/docs/configuration/hostnames)
  - [Auth(-methods)](https://codeshelldev.github.io/secured-signal-api/docs/configuration/auth)
- [Reverse Proxy](https://codeshelldev.github.io/secured-signal-api/docs/reverse-proxy)
- [Integrations](https://codeshelldev.github.io/secured-signal-api/docs/integrations)
- [Contributing](#contributing)
- [Support](#support)
- [Help](#help)
- [License](#license)
- [Legal](#legal)

## Getting Started

> **Prerequisites**: You need Docker and Docker Compose installed.

Get the latest version of the `docker-compose.yaml` file:

```yaml
{{{ #://../../docs/getting-started/examples/docker-compose.yaml }}}
```

And add secure tokens to `api.tokens`. See [API Tokens](#api-tokens).

> [!IMPORTANT]
> Here we'll use `sec-signal-api:8880` as the host,
> but replace it with your actual container/host IP, port, or hostname

## Setup

Before you can send messages via Secured Signal API you must first set up [Signal CLI REST API](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md)

1. **Register** or **link** a Signal account with `signal-cli-rest-api`

2. Deploy `secured-signal-api` with at least one API token

3. Confirm you can send a test message (See [Usage](#usage))

> [!IMPORTANT]
> Run setup directly with Signal CLI REST API.
> Setup requests via Secured Signal API [are blocked by default](#endpoints)

## Usage

Secured Signal API provides 5 ways to authenticate

### Auth

| Method      | Example                                                    |
| :---------- | :--------------------------------------------------------- |
| Bearer Auth | Add `Authorization: Bearer API_TOKEN` to headers           |
| Basic Auth  | Add `Authorization: Basic BASE64_STRING` (`api:API_TOKEN`) |
| Query Auth  | Append `@auth=API_TOKEN` to request URL                    |
| Path Auth   | Prepend request path with `/@auth=API_TOKEN/`              |
| Body Auth   | Set `auth` to `API_TOKEN` in the request body              |

> [!WARNING]
> **Query** and **Path** auth are disabled by default and [must be enabled in the config](https://codeshelldev.github.io/secured-signal-api/docs/configuration/auth)

### Example

To send a message to `+123400002`:

```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer API_TOKEN" -d '{"message": "Hello World!", "recipients": ["+123400002"]}' http://sec-signal-api:8880/v2/send
```

### Advanced

#### Placeholders

If you are not comfortable / don't want to hard-code your number for example and/or recipients in you, may use **placeholders** in your request.

**How to use:**

| Scope                  | Example             | Note             |
| :--------------------- | :------------------ | :--------------- |
| Body                   | `{{@data.key}}`     |                  |
| Header                 | `{{#Content_Type}}` | `-` becomes `_`  |
| [Variable](#variables) | `{{.VAR}}`          | always uppercase |

**Where to use:**

| Scope | Example                                                          |
| :---- | :--------------------------------------------------------------- |
| Body  | `{"number": "{{ .NUMBER }}", "recipients": "{{ .RECIPIENTS }}"}` |
| Query | `http://sec-signal-api:8880/v1/receive/?@number={{.NUMBER}}`     |
| Path  | `http://sec-signal-api:8880/v1/receive/{{.NUMBER}}`              |

You can also combine them:

```json
{
	"content": "{{.NUMBER}} -> {{.RECIPIENTS}}"
}
```

#### URL-to-Body Injection

In some cases you may not be able to access / modify the request body, in that case specify needed values in the request query or path:

`http://sec-signal-api:8880/@key2=value2/?@key=value`

> [!IMPORTANT]
> To differentiate **injection queries** from _regular_ queries, **prefix the key with `@`**.
> Only keys starting with `@` are injected into the request body.

> [!NOTE]
>
> - Supported value types include **strings**, **integers**, **arrays**, and **JSON objects**
> - See [Formatting](https://codeshelldev.github.io/secured-signal-api/docs/usage/formatting) for details on supported structures and syntax

Supported [placeholder types](#placeholders):

| `.` Variables | `@` Body | `#` Headers |
| ------------- | -------- | ----------- |
| ❌            | ✅       | ❌          |

## Configuration

There are multiple ways to configure Secured Signal API, you can optionally use `config.yml` as well as environment variables to override the config.

### Config Files

Config files allow **YAML** formatting and `${ENV}` to get environment variables.

To change the internal config file location set `CONFIG_PATH` in your **environment**. (default: `/config/config.yml`)

This example config shows all the individual settings that can be applied:

```yaml
{{{ #://docs/configuration/examples/config.yml }}}
```

#### Token Configs

You can also override the `config.yml` file for each individual token by adding configs under `TOKENS_PATH` (default: `config/tokens/`)

Here is an example:

```yaml
{{{ #://docs/configuration/examples/token.yml }}}
```

### API Tokens

During authentication Secured Signal API will try to match the given token against the list of tokens inside of the `api.tokens` (or [`api.auth.tokens`](https://codeshelldev.github.io/secured-signal-api/docs/configuration/auth)) attribute.

```yaml
api:
  tokens: [token1, token2, token3]
```

> [!IMPORTANT]
> Using API tokens is highly recommended, but not mandatory.
> Some important security features won't be available (for example the [default blocked endpoints](#endpoints))

> [!NOTE]
> Blocked endpoints can be reactivated by manually configuring them

### Endpoints

Since Secured Signal API is just a proxy you can use all the [Signal CLI REST API](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md) endpoints except for…

| Endpoint              |                    |
| :-------------------- | ------------------ |
| **/v1/configuration** | **/v1/unregister** |
| **/v1/devices**       | **/v1/contacts**   |
| **/v1/register**      | **/v1/accounts**   |
| **/v1/qrcodelink**    |                    |

These endpoints are blocked by default due to security risks.

> [!IMPORTANT]
>
> 1. Matching uses [regex](https://regex101.com)
> 2. On compile error exact match is used instead

> [!WARNING]
> Remember that some symbols have special meanings in regex, a good rule of thumb is:
>
> - If it is a special character, it probably needs to be escaped (`/`) if you are not looking to use regex
> - Otherwise test your pattern on a [regex testing site](https://regex101.com)

You can modify endpoints by configuring `access.endpoints` in your config:

```yaml
settings:
  access:
    endpoints:
      - "!/v1/receive"
      - /v2/send
```

By default adding an endpoint explicitly allows access to it, use `!` to block it instead.

> [!IMPORTANT]
> When using `!` to block you must enclose the endpoint with quotes, like in the example above

| Allow      | Block          | Result                                    |
| ---------- | -------------- | ----------------------------------------- |
| `/v2/send` | —              | **Only** `/v2/send` allowed               |
| —          | `!/v1/receive` | **All** allowed, **except** `/v1/receive` |
| `/v2/send` | `!/v2/.*`      | **Only** `/v2/send` allowed               |

### Variables

Variables can be added under `variables` and can then be referenced in the body, query, or path.
See [Placeholders](#placeholders).

> [!NOTE]
> Variables are always converted into an **uppercase** string.
> Example: `number` ⇒ `NUMBER` in `{{.NUMBER}}`

```yaml
settings:
  message:
    variables:
      number: "+123400001",
      recipients: ["+123400002", "group.id", "user.id"]
```

### Message Templates

To customize the `message` attribute you can use **Message Templates** to build your message by using other body keys and variables.
Use `message.template` to configure:

```yaml
settings:
  message:
    template: |
      Your Message:
      {{@message}}.
      Sent with Secured Signal API.
```

Supported [placeholder types](#placeholders):

| `.` Variables | `@` Body | `#` Headers |
| ------------- | -------- | ----------- |
| ✅            | ✅       | ✅          |

### Field Policies

**Field Policies** allow for blocking or specifically allowing certain fields with set values from being used in the requests body or headers.

Configure them by using `access.fieldPolicies` like so:

```yaml
settings:
  access:
    fieldPolicies:
      "@number":
        - value: "+123400002"
          action: block
        - value: "+123400003"
          action: block
```

Set the wanted action on encounter, available options are `block` and `allow`.

> [!IMPORTANT]
> String fields always try to use
>
> 1. [Regex matching](https://regex101.com)
> 2. On compile error exact match is used as fallback

> [!WARNING]
> Remember that some symbols have special meanings in regex, a good rule of thumb is:
>
> - If it is a special, it probably needs to be escaped (`/`) if you are not looking to use regex
> - Otherwise test your pattern on a [regex testing site](https://regex101.com)

Supported [placeholder types](#placeholders):

| `.` Variables | `@` Body | `#` Headers |
| ------------- | -------- | ----------- |
| ❌            | ✅       | ✅          |

### Field Mappings

To improve compatibility with other services Secured Signal API provides **Field Mappings** and a built-in `message` mapping.

<details>
<summary><strong>Default `message` mapping</strong></summary>

| Field        | Score | Field            | Score |
| ------------ | ----- | ---------------- | ----- |
| msg          | 100   | data.content     | 9     |
| content      | 99    | data.description | 8     |
| description  | 98    | data.text        | 7     |
| text         | 20    | data.summary     | 6     |
| summary      | 15    | data.details     | 5     |
| details      | 14    | body             | 2     |
| data.message | 10    | data             | 1     |

</details>

Secured Signal API will pick the best scoring field (if available) to set the key to the correct value from the request body.

Field Mappings can be added by setting `message.fieldMappings` in your config:

```yaml
settings:
  message:
    fieldMappings:
      "@message":
        [
          { field: "msg", score: 80 },
          { field: "data.message", score: 79 },
          { field: "array[0].message", score: 78 },
        ]
      ".NUMBER": [{ field: "phone_number", score: 100 }]
```

Supported [placeholder types](#placeholders):

| `.` Variables | `@` Body | `#` Headers |
| ------------- | -------- | ----------- |
| ✅            | ✅       | ❌          |

## Contributing

Found a bug? Want to change or add something?
Feel free to open up an [issue](https://github.com/codeshelldev/secured-signal-api/issues) or create a [pull request](https://github.com/codeshelldev/secured-signal-api/pulls)!

## Support

Has this Repo been helpful 👍️ to you? Then consider ⭐️'ing this Project.

:)

## Help

**Are you having problems setting up Secured Signal API?**<br>
No worries check out the [discussions](https://github.com/codeshelldev/secured-signal-api/discussions) tab and ask for help.

**We are all volunteers**, so please be friendly and patient.

## License

This Project is licensed under the [MIT License](./LICENSE).

## Legal

Logo designed by [@CodeShellDev](https://github.com/codeshelldev), All Rights Reserved.

This Project is not affiliated with the Signal Foundation.
