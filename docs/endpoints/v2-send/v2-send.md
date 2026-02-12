---
title: /v2/send
---

# /v2/send

The `/v2/send` endpoint is used to send messages, it is primarily implemented by `signal-cli-rest-api`.

## Payload

In addition to the [standard `/v2/send`](https://bbernhard.github.io/signal-cli-rest-api/#/Messages/post_v2_send) payload,
**Secured Signal API** adds the following fields:

```json
{{{ #://./payload.json }}}
```

### Example

```bash
curl -X POST \
     -H "Authorization: API_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"message":"Happy new year!", "send_at": "1893452400"}' \
      'http://sec-signal-api/v2/send'
```

## Fields

### `send_at`

The `send_at` field is a UNIX timestamp, it tells **Secured Signal API**'s scheduler when to fire the request.

#### Response

| Status Code | Note                  | Fields  |
| :---------: | --------------------- | :-----: |
|   **202**   | Request scheduled     |  `id`   |
|   **400**   | Invalid request       | `error` |
|   **500**   | Internal Server Error |    –    |

```json
{
	"error": "{error}",
	"id": "{uuid}"
}
```

The `id` can be used to check up on the status of the request or to cancel the scheduled request.
