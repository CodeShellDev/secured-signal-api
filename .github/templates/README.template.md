<img align="center" width="1048" height="512" alt="Secure Proxy for Signal REST API" src="https://github.com/CodeShellDev/secured-signal-api/raw/refs/heads/main/logo/landscape" />

<h5 align="center">Secure Proxy for <a href="https://github.com/bbernhard/signal-cli-rest-api">Signal Messenger REST API</a></h5>

## Getting Started

Get the latest version of the `docker-compose.yaml` file:

```yaml
{ { file.docker-compose.yaml } }
```

And add secure Token(s) to `api.tokens`. See [API TOKEN(s)](#api-tokens).

> [!IMPORTANT]
> This Documentation will be using `sec-signal-api:8880` as the service host,
> this **is just for simplicty**, instead use your containers or hosts IP + Port.
> Or a hostname if applicable. See [Reverse Proxy](#reverse-proxy)

### Reverse proxy

Take a look at the [traefik](https://github.com/traefik/traefik) implementation:

```yaml
{ { file.examples/traefik.docker-compose.yaml } }
```

## Setup

Before you can send messages via Secured Signal API you must first setup [Signal rAPI](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md)

To be able to use the API you have to either:

- **register with your Signal Account**

OR

- **link Signal API to an already registered Signal Device**

> [!TIP]
> It is advised to do Setup directly with Signal rAPI
> if you try to Setup with Secured Signal API you will be blocked from doing so. See [Blocked Endpoints](#blocked-endpoints).

## Usage

Secured Signal API provides 3 Ways to Authenticate

### Bearer

To Authenticate add `Authorization: Bearer API_TOKEN` to your request Headers

### Basic Auth

To use Basic Auth as Authorization Method add `Authorization: Basic BASE64_STRING` to your Headers

User is `api` (LOWERCASE)

Formatting for `BASE64_STRING` = `user:API_TOKEN`.

example:

```bash
echo "api:API_TOKEN" | base64
```

=> `YXBpOkFQSV9LRVkK`

### Query Auth

If you are working with a limited Application you may **not** be able to modify Headers or the Request Body
in this case you can use **Query Auth**.

Here is a simple example:

```bash
curl -X POST http://sec-signal-api:8880/v2/send?@authorization=API_TOKEN
```

Notice the `@` infront of `authorization`. See [KeyValue Pair Injection](#keyvalue-pair-injection).

### Example

To send a message to 1234567:

```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer API_TOKEN" -d '{"message": "Hello World!", "recipients": ["1234567"]}' http://sec-signal-api:8880/v2/send
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

## Configuration

There are multiple ways to configure Secured Signal API, you can optionally use `config.yml` aswell as Environment Variables to override the config.

### Config File

Config files allow **YML** formatting and also `${ENV}` to get Environment Variables.

To change the internal config file location set `CONFIG_PATH` in your **Environment** to an absolute path including the filename.extension. (default: `/config/config.yml`)

This example config shows all of the individual settings that can be applied:

```yaml
{ { file.examples/config.yml } }
```

### Environment

Suppose you want to set a new [Placeholder](#placeholders) `NUMBER` in your Environment...

```yaml
environment:
  VARIABLES__NUMBER: "000"
```

This would internally be converted into `variables.number` matching the config formatting.

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
> If you have a string that should not be turned into any other type, then you will need to escape all Type Denotations, `[]` or `{}` (also `-`) with a `\` **Backslash**.
> **Double Backslashes** do exist but you could just leave them out completly.
> An **Odd** number of **Backslashes** **escape** the character in front of them and an **Even** number leave the character **as-is**.

### API Token(s)

During Authentication Secured Signal API will try to match the given Token against the list of Tokens inside of these Variables.

> [!NOTE]
> Both `api.token` and `api.tokens` support multiple Tokens.

```yaml
api:
  token: [token1, token2, token3]
  tokens: [token1, token2, token3]
```

> [!IMPORTANT]
> It is highly recommended use API Tokens

> _What if I just don't?_

Secured Signal API will still work, but important Security Features won't be available
like Blocked Endpoints and any sort of Auth.

> [!NOTE]
> Blocked Endpoints can be reactivated by manually configuring them

### Blocked Endpoints

Because Secured Signal API is just a Proxy you can use all of the [Signal REST API](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md) endpoints except for...

| Endpoint              |
| :-------------------- |
| **/v1/about**         |
| **/v1/configuration** |
| **/v1/devives**       |
| **/v1/register**      |
| **/v1/unregister**    |
| **/v1/qrcodelink**    |
| **/v1/accounts**      |
| **/v1/contacts**      |

These Endpoints are blocked by default due to Security Risks, but can be modified by setting `blockedEndpoints` in your config:

```yaml
blockedEndpoints: [/v1/register, /v1/unregister, /v1/qrcodelink, /v1/contacts]
```

### Variables

Placeholders can be added under `variables` and can then be referenced in the Body, Query or URL.
See [Placeholders](#placeholders).

> [!NOTE]
> Every Placeholder Key will be converted into an Uppercase String.
> Example: `number` becomes `NUMBER` in `{{.NUMBER}}`

```yaml
variables:
  number: "001",
  recipients: [
    "user.id", "000", "001", "group.id"
  ]
```

### Message Aliases

To improve compatibility with other services Secured Signal API provides aliases for the `message` attribute by default:

| Alias       | Score |
| ----------- | ----- |
| msg         | 100   |
| content     | 99    |
| description | 98    |
| text        | 20    |
| body        | 15    |
| summary     | 10    |
| details     | 9     |
| payload     | 2     |
| data        | 1     |

Secured Signal API will pick the best scoring Message Alias (if available) to extract the correct message from the Request Body.

Message Aliases can be added by setting `messageAliases` in your config:

```yaml
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

| Level   |
| ------- |
| `info`  |
| `debug` |
| `warn`  |
| `error` |
| `fatal` |

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
