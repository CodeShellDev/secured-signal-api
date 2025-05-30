# Secured Signal Api

Secured Signal Api acts as a secured proxy for signal-rest-api.

## Installation

Get the latest version of the `docker-compose.yaml` file:

```yaml
{ { file.docker-compose.yaml } }
```

## Usage

To send a message to `number`: `1234567`:

```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer TOKEN" -d '{"message": "Hello World!", "recipients": ["1234567"]}' http://signal-api/v2/send
```

## Contributing

## License

[MIT](https://choosealicense.com/licenses/mit/)
