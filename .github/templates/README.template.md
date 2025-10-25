<img align="center" width="1048" height="512" alt="Secure Proxy for Signal REST API" src="https://github.com/CodeShellDev/secured-signal-api/raw/refs/heads/main/logo/banner.png" />

<h3 align="center">Secure Proxy for <a href="https://github.com/bbernhard/signal-cli-rest-api">Signal Messenger REST API</a></h3>

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
		src="https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fghcr-badge.elias.eu.org%2Fapi%2Fcodeshelldev%2Fsecured-signal-api%2Fsecured-signal-api&query=downloadCount&label=Downloads&color=2344cc11"
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

Check out the official [Documentation](https://codeshelldev.github.io/secured-signal-api) for up-to-date Instructions and additional Content.

- [Getting Started](#getting-started)
- [Setup](#setup)
- [Usage](#usage)
- [Configuration](#configuration)
  - [Endpoints](#endpoints)
  - [Variables](#variables)
  - [Field Mappings](#field-mappings)
  - [Message Templates](#message-templates)
- [Integrations](https://codeshelldev.github.io/secured-signal-api/docs/integrations/compatibility)
- [Contributing](#contributing)
- [Support](#support)
- [Help](#help)
- [License](#license)

## Getting Started

> **Prerequisites**: You need Docker and Docker Compose installed.

Get the latest version of the `docker-compose.yaml` file:

```yaml
{{{ #://docs/getting-started/examples/docker-compose.yaml }}}
```

And add secure Token(s) to `api.tokens`. See [API TOKENs](#api-tokens).

> [!IMPORTANT]
> In this documentation, we use `sec-signal-api:8880` as the host for simplicity.
> Replace it with your actual container/host IP, port, or hostname.

## Setup

Before you can send messages via Secured Signal API you must first set up [Signal rAPI](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md)

1. **Register** or **link** a Signal account with `signal-cli-rest-api`

2. Deploy `secured-signal-api` with at least one API token

3. Confirm you can send a test message (see [Usage](#usage))

> [!TIP]
> Run setup directly with Signal rAPI.
> Setup requests via Secured Signal API are blocked. See [Blocked Endpoints](#endpoints).

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

If you are not comfortable / don't want to hardcode your Number for example and/or Recipients in you, may use **Placeholders** in your Request.

**How to use:**

| Type                   | Example             | Note             |
| :--------------------- | :------------------ | :--------------- |
| Body                   | `{{@data.key}}`     |                  |
| Header                 | `{{#Content_Type}}` | `-` becomes `_`  |
| [Variable](#variables) | `{{.VAR}}`          | always uppercase |

**Where to use:**

| Type  | Example                                                          |
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

#### KeyValue Pair Injection

In some cases you may not be able to access / modify the Request Body, in that case specify needed values in the Request Query:

`http://sec-signal-api:8880/?@key=value`

In order to differentiate Injection Queries and _regular_ Queries
you have to add `@` in front of any KeyValue Pair assignment.

Supported types include **strings**, **ints**, **arrays** and **json dictionaries**. See [Formatting](https://codeshelldev.github.io/secured-signal-api/docs/usage/formatting).

## Configuration

There are multiple ways to configure Secured Signal API, you can optionally use `config.yml` aswell as Environment Variables to override the config.

### Config Files

Config files allow **YAML** formatting and also `${ENV}` to get Environment Variables.

To change the internal config file location set `CONFIG_PATH` in your **Environment** to an absolute path including the filename.extension. (default: `/config/config.yml`)

This example config shows all of the individual settings that can be applied:

```yaml
{{{ #://docs/configuration/examples/config.yml }}}
```

#### Token Configs

You can also override the `config.yml` file for each individual token by adding configs under `TOKENS_PATH` (default: `config/tokens/`)

This way you can permission tokens by further restricting or adding [Endpoints](#endpoints), [Placeholders](#variables), etc.

Here is an example:

```yaml
{{{ #://docs/configuration/examples/token.yml }}}
```

### Templating

Secured Signal API uses Golang's [Standard Templating Library](https://pkg.go.dev/text/template).
This means that any valid Go template string will also work in Secured Signal API.

Go's templating library is used in the following features:

- [Message Templates](#message-templates)
- [Placeholders](#placeholders)

This makes advanced [Message Templates](#message-templates) like this one possible:

```yaml
{{{ #://docs/configuration/examples/message-template.yml }}}
```

### API Tokens

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

You can modify endpoints by configuring `access.endpoints` in your config:

```yaml
settings:
  access:
    endpoints:
      - !/v1/register
      - !/v1/unregister
      - !/v1/qrcodelink
      - !/v1/contacts
      - /v2/send
```

By default adding an endpoint explictly allows access to it, use `!` to block it instead.

| Config (Allow) | (Block)        |   Result   |     |                   |     |
| :------------- | :------------- | :--------: | --- | :---------------: | --- |
| `/v2/send`     | `unset`        |  **all**   | 🛑  |  **`/v2/send`**   | ✅  |
| `unset`        | `!/v1/receive` |  **all**   | ✅  | **`/v1/receive`** | 🛑  |
| `/v2`          | `!/v2/send`    | **`/v2*`** | 🛑  |  **`/v2/send`**   | ✅  |

### Variables

Placeholders can be added under `variables` and can then be referenced in the Body, Query or URL.
See [Placeholders](#placeholders).

> [!NOTE]
> Every Placeholder Key will be converted into an Uppercase String.
> Example: `number` becomes `NUMBER` in `{{.NUMBER}}`

```yaml
settings:
  message:
    variables:
      number: "+123400001",
      recipients: ["+123400002", "group.id", "user.id"]
```

### Message Templates

To customize the `message` attribute you can use **Message Templates** to build your message by using other Body Keys and Variables.
Use `message.template` to configure:

```yaml
settings:
  message:
    template: |
      Your Message:
      {{@message}}.
      Sent with Secured Signal API.
```

Message Templates support [Standard Golang Templating](#templating).
Use `@data.key` to reference Body Keys, `#Content_Type` for Headers and `.KEY` for Variables.

### Field Mappings

To improve compatibility with other services Secured Signal API provides **Field Mappings** and a built-in `message` Mapping.

<details>
<summary><strong>Default `message` Mapping</strong></summary>

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

Secured Signal API will pick the best scoring Field (if available) to set the Key to the correct Value from the Request Body.

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

Use `@` for mapping to Body Keys and `.` for mapping to Variables.

## Contributing

Found a bug? Want to change or add something?
Feel free to open up an [Issue](https://github.com/codeshelldev/secured-signal-api/issues) or create a [Pull Request](https://github.com/codeshelldev/secured-signal-api/pulls)!

## Support

Has this Repo been helpful 👍️ to you? Then consider ⭐️'ing this Project.

:)

## Help

**Are you having Problems setting up Secured Signal API?**<br>
No worries check out the [Discussions](https://github.com/codeshelldev/secured-signal-api/discussions) Tab and ask for help.

**We are all Volunteers**, so please be friendly and patient.

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Legal

Logo designed by [@CodeShellDev](https://github.com/codeshelldev), All Rights Reserved.

This Project is not affiliated with the Signal Foundation.
