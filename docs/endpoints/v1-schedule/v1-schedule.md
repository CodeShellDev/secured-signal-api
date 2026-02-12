---
title: /v1/schedule
---

# /v1/schedule

The `/v1/schedule` is a custom endpoint only implemented by **Secured Signal API**.

## DELETE

### Example

```bash
curl -X DELETE -H "Authorization: API_TOKEN" http://sec-signal-api/v1/schedule/{uuid}
```

### Response

| Status Code | Note             | Fields  |
| :---------: | ---------------- | :-----: |
|   **204**   | Request canceled |    –    |
|   **400**   | Invalid id       | `error` |

## GET

### Example

```bash
curl -H "Authorization: API_TOKEN" http://sec-signal-api/v1/schedule/{uuid}
```

> [!NOTE]
> The `{uuid}` is returned by the [`/v2/send` endpoint](./v2-send#send_at) in the response body as `id`

### Response

| Status Code | Note             | Fields  |
| :---------: | ---------------- | :-----: |
|   **200**   | Successful fetch |  `...`  |
|   **400**   | Invalid id       | `error` |

#### Pending / Queued

```json
{{{ #://./response/get_pending.json }}}
```

#### Failed

```json
{{{ #://./response/get_failed.json }}}
```

#### Done

```json
{{{ #://./response/get_done.json }}}
```

The `response_body`, `response_headers` and `response_status_code` are nullable.
An **empty response body** or **empty response headers** are **`null`** and not `{ }` or `[ ]` (_empty_).

> [!IMPORTANT]
> After a finished request has been fetched via this endpoint the entry is automatically deleted from the database
