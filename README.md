<img align="center" width="1048" height="512" alt="Secure Proxy for Signal CLI REST API" src="https://github.com/CodeShellDev/secured-signal-api/raw/refs/heads/docs/static/img/banner.png" />

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

Check out the [**Official Documentation**](https://codeshelldev.github.io/secured-signal-api) for up-to-date instructions and additional content!

- [Getting Started](#getting-started)
- [Setup](#setup)
- [Usage](#usage)
- [Configuration](#configuration)
  - [Endpoints](#endpoints)
  - [Variables](#variables)
  - [Field Policies](#field-policies)
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
    restart: unless-stopped
    networks:
      backend:
        aliases:
          - secured-signal-api

networks:
  backend:
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

Secured Signal API provides 3 ways to authenticate

### Auth

| Method      | Example                                                    |
| :---------- | :--------------------------------------------------------- |
| Bearer Auth | Add `Authorization: Bearer API_TOKEN` to headers           |
| Basic Auth  | Add `Authorization: Basic BASE64_STRING` (`api:API_TOKEN`) |
| Query Auth  | Append `@auth=API_TOKEN` to request URL                    |

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

#### Query-to-Body Injection

In some cases you may not be able to access / modify the request body, in that case specify needed values in the request query:

`http://sec-signal-api:8880/?@key=value`

> [!IMPORTANT]
> To differentiate **injection queries** from _regular_ queries, **prefix the key with `@`**.
> Only keys starting with `@` are injected into the request body

> [!NOTE]
>
> - Supported value types include **strings**, **integers**, **arrays**, and **JSON objects**
> - See [Formatting](https://codeshelldev.github.io/secured-signal-api/docs/usage/formatting) for details on supported structures and syntax

## Configuration

There are multiple ways to configure Secured Signal API, you can optionally use `config.yml` as well as environment variables to override the config.

### Config Files

Config files allow **YAML** formatting and `${ENV}` to get environment variables.

To change the internal config file location set `CONFIG_PATH` in your **environment**. (default: `/config/config.yml`)

This example config shows all the individual settings that can be applied:

```yaml
# Example Config (all configurations shown)
service:
  port: 8880

api:
  url: http://signal-api:8080
  tokens: [token1, token2]

logLevel: info

settings:
  message:
    template: |
      You've got a Notification:
      {{@message}} 
      At {{@data.timestamp}} on {{@data.date}}.
      Send using {{.NUMBER}}.

    variables:
      number: "+123400001"
      recipients: ["+123400002", "group.id", "user.id"]

    fieldMappings:
      "@message": [{ field: "msg", score: 100 }]

  access:
    endpoints:
      - "!/v1/about"
      - /v2/send

    fieldPolicies:
      "@number": {
        value: "+123400003",
        action: block
      }
```

#### Token Configs

You can also override the `config.yml` file for each individual token by adding configs under `TOKENS_PATH` (default: `config/tokens/`)

Here is an example:

```yaml
api:
  tokens: [LOOOONG_STRING]

settings:
  message:
    fieldMappings: # Disable mappings
    variables: # Disable variable placeholders

  access:
    endpoints: # Disable sending
      - "!/v2/send"
```

### Templating

Secured Signal API uses Go's [standard templating library](https://pkg.go.dev/text/template).
This means that any valid Go template string will also work in Secured Signal API.

Go's templating library is used in the following features:

- [Message Templates](#message-templates)
- [Query-to-Body Injection](#query-to-body-injection)
- [Placeholders](#placeholders)

This makes advanced [Message Templates](#message-templates) like this one possible:

```yaml
settings:
  message:
    template: |
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
      ------------------------------------
      Content-Type: {{ #Content_Type }}
      Redacted Auth Header: {{ #Authorization }}
```

### API Tokens

During authentication Secured Signal API will try to match the given token against the list of tokens inside of the `api.tokens` attribute.

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

> [!NOTE]
> Matching uses [glob-like patterns](https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html):
>
> - `*` matches any sequence of characters
> - `?` matches a single character
> - `[abc]` matches one of the characters in the brackets

You can modify endpoints by configuring `access.endpoints` in your config:

```yaml
settings:
  access:
    endpoints:
      - "!/v1/register"
      - "!/v1/unregister"
      - "!/v1/qrcodelink"
      - "!/v1/contacts"
      - /v2/send
```

By default adding an endpoint explicitly allows access to it, use `!` to block it instead.

> [!IMPORTANT]
> When using `!` to block you must enclose the endpoint with quotes, like in the example above

| Config (Allow) | (Block)        |   Result   |     |                   |     |
| :------------- | :------------- | :--------: | --- | :---------------: | --- |
| `/v2/send`     | `unset`        |  **all**   | 🛑  |  **`/v2/send`**   | ✅  |
| `unset`        | `!/v1/receive` |  **all**   | ✅  | **`/v1/receive`** | 🛑  |
| `!/v2*`        | `/v2/send`     | **`/v2*`** | 🛑  |  **`/v2/send`**   | ✅  |

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
      "@number": { value: "+123400002", action: block }
```

Set the wanted action on encounter, available options are `block` and `allow`.

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
