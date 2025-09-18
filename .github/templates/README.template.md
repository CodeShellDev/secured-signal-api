<img align="center" width="1048" height="512" alt="Secure Proxy for Signal REST API" src="https://github.com/CodeShellDev/secured-signal-api/raw/refs/heads/main/logo/landscape" />

<h3 align="center">Secure Proxy for <a href="https://github.com/bbernhard/signal-cli-rest-api">Signal Messenger REST API</a></h3>

<p align="center">
adding token-based authentication,
endpoint restrictions, placeholders, and flexible configuration.
</p>

<p align="center">
🔒 Secure · ⭐️ Configurable · 🚀 Easy to Deploy with Docker
</p>

<div align="center">
  <a href="https://github.com/codeshelldev/secured-signal-api/releases">
    <img src="https://img.shields.io/github/v/release/codeshelldev/secured-signal-api?sort=semver&logo=github" alt="GitHub release">
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/pkgs/container/secured-signal-api">
    <img src="https://ghcr-badge.egpl.dev/codeshelldev/secured-signal-api/size?color=%2344cc11&tag=latest&label=image+size&trim=" alt="Docker image size">
  </a>
  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT">
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/stargazers">
    <img src="https://img.shields.io/github/stars/codeshelldev/secured-signal-api?style=flat&logo=github" alt="GitHub stars">
  </a>
</div>

## Contents

- [Getting Started](#getting-started)
- [Setup](#setup)
- [Usage](#usage)
- [Best Practices](#security-best-practices)
- [Configuration](#configuration)
  - [Endpoints](#endpoints)
  - [Variables](#variables)
- [Contributing](#contributing)
- [Support](#support)
- [License](#license)

## Getting Started

> **Prerequisites**: You need Docker and Docker Compose installed.

Get the latest version of the `docker-compose.yaml` file:

```yaml
{ { file.docker-compose.yaml } }
```

And add secure Token(s) to `api.tokens`. See [API TOKENs](#api-tokens).

> [!IMPORTANT]
> In this documentation, we use `sec-signal-api:8880` as the host for simplicity.
> Replace it with your actual container/host IP, port, or hostname.

### Reverse Proxy

Take a look at the [traefik](https://github.com/traefik/traefik) implementation:

```yaml
{ { file.examples/traefik.docker-compose.yaml } }
```

## Setup

Before you can send messages via Secured Signal API you must first set up [Signal rAPI](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md)

1. **Register** or **link** a Signal account with `signal-cli-rest-api`

2. Deploy `secured-signal-api` with at least one API token

3. Confirm you can send a test message (see [Usage](#usage))

> [!TIP]
> Run setup directly with Signal rAPI.
> Setup requests via Secured Signal API are blocked. See [Blocked Endpoints](#blocked-endpoints).

## Usage

Secured Signal API provides 3 Ways to Authenticate

### Auth

| Method      | Example                                                    |
| :---------- | :--------------------------------------------------------- |
| Bearer Auth | Add `Authorization: Bearer API_TOKEN` to headers           |
| Basic Auth  | Add `Authorization: Basic BASE64_STRING` (`api:API_TOKEN`) |
| Query Auth  | Append `@authorization=API_TOKEN` to request URL           |

### Example

To send a message to `+123400002`:

```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer API_TOKEN" -d '{"message": "Hello World!", "recipients": ["+123400002"]}' http://sec-signal-api:8880/v2/send
```

### Advanced

#### Placeholders

If you are not comfortable / don't want to hardcode your Number for example and/or Recipients in you, may use **Placeholders** in your Request. See [Custom Variables](#variables).

These Placeholders can be used in the Request Query or the Body of a Request like so:

**Body**

```json
{
	"number": "{{ .NUMBER }}",
	"recipients": "{{ .RECIPIENTS }}"
}
```

**Query**

```
http://sec-signal-api:8880/v1/receive/?@number={{.NUMBER}}
```

**Path**

```
http://sec-signal-api:8880/v1/receive/{{.NUMBER}}
```

#### KeyValue Pair Injection

In some cases you may not be able to access / modify the Request Body, in that case specify needed values in the Request Query:

`http://sec-signal-api:8880/?@key=value`

In order to differentiate Injection Queries and _regular_ Queries
you have to add `@` in front of any KeyValue Pair assignment.

Supported types include **strings**, **ints** and **arrays**. See [Formatting](#string-to-type).

## Security: Best Practices

- Always use API tokens in production
- Run behind a TLS-enabled [Reverse Proxy](#reverse-proxy) (Traefik, Nginx, Caddy)
- Be cautious when overriding Blocked Endpoints
- Use per-token overrides to enforce least privilege

## Configuration

There are multiple ways to configure Secured Signal API, you can optionally use `config.yml` aswell as Environment Variables to override the config.

### Config Files

Config files allow **YML** formatting and also `${ENV}` to get Environment Variables.

To change the internal config file location set `CONFIG_PATH` in your **Environment** to an absolute path including the filename.extension. (default: `/config/config.yml`)

This example config shows all of the individual settings that can be applied:

```yaml
{ { file.examples/config.yml } }
```

#### Token Configs

You can also override the `config.yml` file for each individual token by adding configs under `TOKENS_PATH` (default: `config/tokens/`)

This way you can permission tokens by further restricting or adding [Endpoints](#blocked-endpoints), [Placeholders](#variables), etc.

Here is an example:

```yaml
{ { file.examples/token.yml } }
```

### Environment

Suppose you want to set a new [Placeholder](#placeholders) `NUMBER` in your Environment...

```yaml
environment:
  SETTINGS__VARIABLES__NUMBER: "+123400001"
```

This would internally be converted into `settings.variables.number` matching the config formatting.

> [!IMPORTANT]
> Underscores `_` are removed during Conversion, Double Underscores `__` on the other hand convert the Variable into a nested Object (`__` replaced by `.`)

### String To Type

> [!TIP]
> This formatting applies to almost every situation where the only (allowed) Input Type is a string and other Output Types are needed.

If you are using Environment Variables as an example you won't be able to specify an Array or a Dictionary of items, in that case you can provide a specifically formatted string which will be translated into the correct type...

| type       | example           |
| :--------- | :---------------- |
| string     | abc               |
| string     | +123              |
| int        | 123               |
| int        | -123              |
| json       | {"a":"b","c":"d"} |
| array(int) | [1,2,3]           |
| array(str) | [a,b,c]           |

> [!NOTE]
> If you have a string that should not be turned into any other type, then you will need to escape all Type Denotations, `[]` or `{}` (also `-`) with a `\` **Backslash** (or Double Backslash).
> An **Odd** number of **Backslashes** **escape** the character in front of them and an **Even** number leave the character **as-is**.

### API Token(s)

During Authentication Secured Signal API will try to match the given Token against the list of Tokens inside of these Variables.

```yaml
api:
  tokens: [token1, token2, token3]
```

> [!IMPORTANT]
> Using API Tokens is highly recommended, but not mandatory.
> Some important Security Features won't be available (like default Blocked Endpoints).

> [!NOTE]
> Blocked Endpoints can be reactivated by manually configuring them

### Endpoints

Since Secured Signal API is just a Proxy you can use all of the [Signal REST API](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md) endpoints except for...

| Endpoint              |                    |
| :-------------------- | ------------------ |
| **/v1/about**         | **/v1/unregister** |
| **/v1/configuration** | **/v1/qrcodelink** |
| **/v1/devices**       | **/v1/contacts**   |
| **/v1/register**      | **/v1/accounts**   |

These Endpoints are blocked by default due to Security Risks.

> [!NOTE]
> Matching works by checking if the requested Endpoints starts with a Blocked or an Allowed Endpoint

You can modify Blocked Endpoints by configuring `blockedEndpoints` in your config:

```yaml
settings:
  blockedEndpoints: [/v1/register, /v1/unregister, /v1/qrcodelink, /v1/contacts]
```

You can also override Blocked Endpoints by adding Allowed Endpoints to `allowedEndpoints`.

```yaml
settings:
  allowedEndpoints: [/v2/send]
```

| Config (Allow)                   | (Block)                             |   Result   |     |                   |     |
| :------------------------------- | :---------------------------------- | :--------: | --- | :---------------: | --- |
| `allowedEndpoints: ["/v2/send"]` | `unset`                             |  **all**   | 🛑  |  **`/v2/send`**   | ✅  |
| `unset`                          | `blockedEndpoints: ["/v1/receive"]` |  **all**   | ✅  | **`/v1/receive`** | 🛑  |
| `blockedEndpoints: ["/v2"]`      | `allowedEndpoints: ["/v2/send"]`    | **`/v2*`** | 🛑  |  **`/v2/send`**   | ✅  |

### Variables

Placeholders can be added under `variables` and can then be referenced in the Body, Query or URL.
See [Placeholders](#placeholders).

> [!NOTE]
> Every Placeholder Key will be converted into an Uppercase String.
> Example: `number` becomes `NUMBER` in `{{.NUMBER}}`

```yaml
settings:
  variables:
    number: "+123400001",
    recipients: ["+123400002", "group.id", "user.id"]
```

### Message Aliases

To improve compatibility with other services Secured Signal API provides **Message Aliases** for the `message` attribute.

<details>
<summary><strong>Default Message Aliases</strong></summary>

| Alias        | Score | Alias            | Score |
| ------------ | ----- | ---------------- | ----- |
| msg          | 100   | data.content     | 9     |
| content      | 99    | data.description | 8     |
| description  | 98    | data.text        | 7     |
| text         | 20    | data.summary     | 6     |
| summary      | 15    | data.details     | 5     |
| details      | 14    | body             | 2     |
| data.message | 10    | data             | 1     |

</details>

Secured Signal API will pick the best scoring Message Alias (if available) to extract the correct message from the Request Body.

Message Aliases can be added by setting `messageAliases` in your config:

```yaml
settings:
  messageAliases:
    [
      { alias: "msg", score: 80 },
      { alias: "data.message", score: 79 },
      { alias: "array[0].message", score: 78 },
    ]
```

### Port

To change the Port which Secured Signal API uses, you need to set `server.port` in your config. (default: `8880`)

### Log Level

To change the Log Level set `logLevel` to: (default: `info`)

<details>
<summary>Log Levels</summary>

| Level   |
| ------- |
| `info`  |
| `debug` |
| `warn`  |
| `error` |
| `fatal` |
| `dev`   |

</details>

## Contributing

Found a bug? Want to change or add something?
Feel free to open up an issue or create a Pull Request!

## Support

Has this Repo been helpful 👍️ to you? Then consider ⭐️'ing this Project.

:)

## License

[MIT](https://choosealicense.com/licenses/mit/)

### Legal

Logo designed by [@CodeShellDev](https://github.com/codeshelldev), All Rights Reserved.
This Project is not affiliated with the Signal Foundation.
