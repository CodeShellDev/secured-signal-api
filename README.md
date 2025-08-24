# Secured Signal API

Secured Signal API acts as a secure proxy for [Signal rAPI](https://github.com/bbernhard/signal-cli-rest-api).

## Installation

Get the latest version of the `docker-compose.yaml` file:

And add secure Token(s) to `API_TOKEN` / `API_TOKENS`. See [API_TOKEN(s)](#api-tokens)

> [!IMPORTANT]
> This Documentation will be using `sec-signal-api:8880` as the service host,
> this **won't work**, instead use your containers IP + Port.
> Or a hostname if applicable. See [Reverse Proxy](#reverse-proxy)

```yaml
services:
  signal-api:
    image: bbernhard/signal-cli-rest-api:latest
    container_name: signal-api
    environment:
      - MODE=normal
    volumes:
      - ./data:/home/.local/share/signal-cli
    networks:
      backend:
        aliases:
          - signal-api
    restart: unless-stopped

  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api:latest
    container_name: secured-signal
    networks:
      backend:
        aliases:
          - secured-signal-api
    environment:
      SIGNAL_API_URL: http://signal-api:8080
      DEFAULT_RECIPIENTS: '[ "000", "001", "002" ]'
      NUMBER: 123456789
      API_TOKEN: LOOOOOONG_STRING
    ports:
      - "8880:8880"
    restart: unless-stopped

networks:
  backend:
```

### Reverse proxy

Take a look at the [traefik](https://github.com/traefik/traefik) implementation:

```yaml
services:
  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api:latest
    container_name: secured-signal
    networks:
      proxy:
      backend:
        aliases:
          - secured-signal-api
    environment:
      SIGNAL_API_URL: http://signal-api:8080
      DEFAULT_RECIPIENTS: '[ "000", "001", "002" ]'
      NUMBER: 123456789
      API_TOKEN: LOOOOOONG_STRING
    labels:
      - traefik.enable=true
      - traefik.http.routers.signal-api.rule=Host(`signal-api.mydomain.com`)
      - traefik.http.routers.signal-api.entrypoints=websecure
      - traefik.http.routers.signal-api.tls=true
      - traefik.http.routers.signal-api.tls.certresolver=cloudflare
      - traefik.http.routers.signal-api.service=signal-api-svc
      - traefik.http.services.signal-api-svc.loadbalancer.server.port=8880
      - traefik.docker.network=proxy
    restart: unless-stopped

networks:
  backend:
  proxy:
    external: true
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

Notice the `@` infront of `authorization`. See [Formatting](#format).

### Example

To send a message to 1234567:

```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer API_TOKEN" -d '{"message": "Hello World!", "recipients": ["1234567"]}' http://sec-signal-api:8880/v2/send
```

### Advanced

#### Placeholders

If you are not comfortable / don't want to hardcode your Number and/or Recipients in you may use **Placeholders** in your Request.

Built-in Placeholders: `{{ .NUMBER }}` and `{{ .RECIPIENTS }}`

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

Supported types include **strings**, **ints** and **arrays**

`http://sec-signal-api:8880/?@key=value`

| type       | example |
| :--------- | :------ |
| string     | abc     |
| int        | 123     |
| array      | [1,2,3] |
| array(int) | 1,2,3   |
| array(str) | a,b,c   |

##### Format

In order to differentiate Injection Queries and _regular_ Queries
you have to add `@` in front of any KeyValue Pair assignment.

### Environment Variables

#### API Token(s)

Both `API_TOKEN` and `API_TOKENS` support multiple Tokens seperated by a `,` **Comma**.
During Authentikcation Secured Signal API will try to match the given Token against the list of Tokens inside of these Variables.

```yaml
environment:
  API_TOKEN: "token1, token2, token3"
  API_TOKENS: "token1, token2, token3"
```

> [!IMPORTANT]
> It is highly recommended to set this Environment Variable

> _What if I just don't?_

Secured Signal API will still work, but important Security Features won't be available
like Blocked Endpoints and any sort of Auth.

> [!NOTE]
> Blocked Endpoints can be reactivated by manually setting them in the Environment

#### Blocked Endpoints

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

These Endpoints are blocked by default due to Security Risks, but can be modified by setting `BLOCKED_ENDPOINTS` to a valid json array string

```yaml
environment:
  BLOCKED_ENDPOINTS: '[ "/v1/register","/v1/unregister","/v1/qrcodelink","/v1/contacts" ]'
```

#### Variables

By default Secured Signal API provides the following Placeholders:

- **NUMBER** = _ENV_: `NUMBER`
- **RECIPIENTS** = _ENV_: `RECIPIENTS`

#### Customization

Placeholders can be added by setting `VARIABLES` inside your Environment.

```yaml
environment:
  VARIABLES: ' "NUMBER2": "002", "GROUP_CHAT_1": [ "user.id", "000", "001", "group.id" ] '
```

#### Recipients

Set this Environment Variable to automatically provide default Recipients:

```yaml
environment:
  RECIPIENTS: ' [ "user.id", "000", "001", "group.id" ] '
```

example:

```json
{
	"recipients": "{{.RECIPIENTS}}"
}
```

#### Message Aliases

To improve compatibility with other services Secured Signal API provides aliases for the `message` attribute by default:

| Alias       | Priority |
| ----------- | -------- |
| msg         | 100      |
| content     | 99       |
| description | 98       |
| text        | 20       |
| body        | 15       |
| summary     | 10       |
| details     | 9        |
| payload     | 2        |
| data        | 1        |

Secured Signal API will use the highest priority Message Alias to extract the correct message from the Request Body.

Message Aliases can be added by setting `MESSAGE_ALIASES`:

```yaml
environment:
  MESSAGE_ALIASES: ' [{ "alias": "note", "priority": 4 }, { "alias": "test", "priority": 3 }] '
```

## Contributing

Found a bug? Want to change or add something?
Feel free to open up an issue or create a Pull Request!

_This is a small project so don't expect any huge changes in the future_

## Support

Has this Repo been helpful 👍️ to you? Then consider ⭐️'ing this Project.

:)

## License

[MIT](https://choosealicense.com/licenses/mit/)
