# Secured Signal Api

Secured Signal Api acts as a secure proxy for signal-rest-api.

## Installation

Get the latest version of the `docker-compose.yaml` file:

And set `API_TOKEN` to a long secure string

```yaml
---
services:
  signal-api:
    image: bbernhard/signal-cli-rest-api
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
    image: ghcr.io/codeshelldev/secured-signal-api
    container_name: secured-signal
    networks:
      backend:
        aliases:
          - secured-signal-api
    environment:
      SIGNAL_API_URL: http://signal-api:8080
      DEFAULT_RECIPIENTS: '[ "000", "001", "002" ]'
      SENDER: 123456789
      API_TOKEN: LOOOOOONG_STRING
    ports:
      - "8880:8880"
    restart: unless-stopped

networks:
  backend:
```

### Reverse proxy

Take a look at traefik implementation:

```yaml
services:
  # ...
  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api
    container_name: secured-signal
    networks:
      proxy:
      backend:
        aliases:
          - secured-signal-api
    environment:
      SIGNAL_API_URL: http://signal-api:8080
      DEFAULT_RECIPIENTS: '[ "000", "001", "002" ]'
      SENDER: 123456789
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

Before you can send messages via `secured-signal-api` you must first setup [`signal-api`](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md),

to send messages you have to either:

- **register a Signal Account**

OR

- **link Signal API to an already registered Signal Device**

## Usage

Secured Signal API implements 3 Ways to Authenticate

### Bearer

To Authenticate with `secured-signal-api` add `Authorization: Bearer TOKEN` to your request Headers

### Basic Auth

To use Basic Auth as Authorization Method add `Authorization: Basic base64{user:pw}` to your Headers

### Query Auth

If you are working with a limited Application you may **not** be able to modify Headers or the Request Body
in this case you should use **Query Auth**.

Here is a simple example:

```bash
curl -X POST http://signal-api:8880/v2/send?@authorization=TOKEN
```

### Example

To send a message to `number`: `1234567`:

```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer TOKEN" -d '{"message": "Hello World!", "recipients": ["1234567"]}' http://signal-api:8880/v2/send
```

### Advanced

#### Placeholders

If you are not comfortable with hardcoding your Number and/or Recipients in you may use **Placeholders** in your request like:

`{{ .NUMBER }}` or `{{ .RECIPIENTS }}`

These _Placeholders_ can be used in the Query or the Body of a Request like so:

**Body**

```json
{
	"number": "{{ .NUMBER }}",
	"recipients": "{{ .RECIPIENTS }}"
}
```

**Query**

```
http://.../?@number={{.NUMBER}}
```

**Path**

```
http://signal-api:8880/v1/receive/{{.NUMBER}}
```

#### KeyValue Pair Injection

In some cases you may not be able to access / modify the Request Body, if that is the case specify needed values in the Requests Query:

```
http://signal-api:8880/?@key=value
```

**Format**
In order to differentiate Injection Queries and _regular_ Queries
you have to add `@` in front of any KeyValue Pair assignment

### Environment Variables

#### API Token

> [!IMPORTANT]
> It is highly recommended to set this Environment Variable to a long secure string

_What if I just don't?_

Well Secured Signal API will still work, but important Security Features won't be available
like Blocked Endpoints and anyone with access to your Docker Container will be able to send Messages in your Name

> [!NOTE]
> Blocked Endpoints can be reactivated by manually setting them in the environment

#### Blocked Endpoints

Because Secured Signal API is just a secure Proxy you can use all of the [Signal REST API](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md) endpoints with an Exception of:

- **/v1/about**

- **/v1/configuration**

- **/v1/devices**

- **/v1/register**

- **/v1/unregister**

- **/v1/qrcodelink**

- **/v1/accounts**

- **/v1/contacts**

These Endpoints are blocked by default to Security Risks, but can be modified by setting `BLOCKED_ENDPOINTS` in the environment variable to a valid json array string

```yaml
environment:
  BLOCKED_ENDPOINTS: '[ "/v1/register","/v1/unregister","/v1/qrcodelink","/v1/contacts" ]'
```

#### Variables

By default Secured Signal API provides the following **Placeholders**:

- **NUMBER** = _ENV_: `SENDER`
- **RECIPIENTS** = _ENV_: `DEFAULT_RECIPIENTS`

If you are ever missing any **Placeholder** (that isn't built-in) you can add as many as you like to `VARIABLES` inside your environment

```yaml
environment:
  VARIABLES: ' "NUMBER2": "002", "GROUP_CHAT_1": [ "user.id", "000", "001", "group.id" ] '
```

#### Default Recipients

Set this environment variable to automatically provide default Recipients:

```yaml
environment:
  DEFAULT_RECIPIENTS: ' [ "user.id", "000", "001", "group.id" ] '
```

## Contributing

Found a bug? Want to change or add something?
Feel free to open up an issue or create a Pull Request!

_This is a small project so don't expect any huge changes in the future_

## License

[MIT](https://choosealicense.com/licenses/mit/)
