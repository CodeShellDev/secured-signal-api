---
sidebar_position: 3
title: Setup
---

# Setup

> [!WARNING]
> These instructions are for **personal or educational use only**. Using multiple accounts, automated messaging, or any activity that violates **Signal's Terms of Service** may result in **account suspension** or **legal action**. We **do not** endorse **spam or fraudulent activity**!
> Furthermore we are **not in any way affiliated** with the **Signal Foundation**.

To use **Secured Signal API** for the first time you will need to set up and link your **Signal Account**.
In this section we'll taking you quickly through what's needed.

> [!TIP]
> Run setup directly with Signal CLI REST API.
> Setup requests via Secured Signal API will be blocked by default. See [Blocked Endpoints](../configuration/endpoints).

## Register

Before sending messages (etc.) via **Secured Signal API** you must first set up Signal CLI REST API.
Here we'll be registering a new Signal Account.

### SMS Verification

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    'http://signal-cli-rest-api:8080/v1/register/<number>'
```

Example:

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    'http://signal-cli-rest-api:8080/v1/register/+431212131491291'
```

### Voice Verification

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"use_voice": true}' \
    'http://signal-cli-rest-api:8080/v1/register/<number>'
```

Example:

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"use_voice": true}' \
    'http://signal-cli-rest-api:8080/v1/register/+431212131491291'
```

## Link

If you don't want to register a new Account then you can instead link a device.

```bash
curl -X GET \
    -H "Content-Type: application/json" \
    'http://signal-cli-rest-api:8080/v1/qrcodelink?device_name=<device name>'
```

This will show you a QR-Code which you will be able to use for linking.

## Troubleshooting

If you encounter any issues in the steps above look at the [examples](https://github.com/bbernhard/signal-cli-rest-api/blob/master/doc/EXAMPLES.md) provided by [@bbernhard](https://github.com/bbernhard/signal-cli-rest-api)
