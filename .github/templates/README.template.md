# Secured Signal Api

Secured Signal Api acts as a secured proxy for signal-rest-api.

## Installation

Get the latest version of the `docker-compose.yaml` file:

```yaml
{ { file.docker-compose.yaml } }
```

### Reverse proxy

Take a look at traefik implementation:

```yaml
{ { file.examples/traefik.docker-compose.yaml } }
```

## Setup

Before you can send messages via `secured-signal-api` you must first setup [`signal-api`](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md),

to send messages you have to either:

- register a Signal Account

OR

- link Signal Api to a already registered Signal Device

## Usage

To send a message to `number`: `1234567`:

```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer TOKEN" -d '{"message": "Hello World!", "recipients": ["1234567"]}' http://signal-api:8880/v2/send
```

### Configuration

Because `secured-signal-api` is just a secure proxy you can use all of the [Signal REST Api](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md) endpoints with an Exception of:

```python
DEFAULT_BLOCKED_ENDPOINTS = [
    "/v1/about",
    "/v1/configuration",
    "/v1/devices",
    "/v1/register",
    "/v1/unregister",
    "/v1/qrcodelink",
    "/v1/accounts",
    "/v1/contacts"
]
```

Which are blocked by default to increase Security, but you these can be modified by setting the `BLOCKED_ENDPOINTS` environment variable as a valid json array

```yaml
environment:
  BLOCKED_ENDPOINTS: '[ "/v1/register","/v1/unregister","/v1/qrcodelink","/v1/contacts" ]'
```

## Contributing

Found a bug? Want to change or add something?
Feel free to open up an issue or create a Pull Request!

_This is a small project so don't expect any huge changes in the future_

## License

[MIT](https://choosealicense.com/licenses/mit/)
