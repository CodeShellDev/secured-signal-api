---
title: Log Levels
---

# Log Levels

Log levels help to filter or explicitly allow certain information from the logs.

To change the log level set `service.logLevel` to one of the following levels… (default: `info`)

**Levels:**

- `info`
- `debug` (verbose)
- `warn` (**only** warnings and errors)
- `error` (**only** errors)
- `fatal` (**only** fatal errors)

> [!CAUTION]
> The log level `dev` **can leak data in the logs**
> and must only be used for testing during development

## Per-Token Logger

Each token config can define its **own log level**, independent of the global logger.

This allows fine-grained control, for example:

- verbose logging for one integration
- minimal logging for another

If `service.logLevel` is not set, the global log level is used.

### Logger Naming

Each token logger is automatically named using the [`name` attribute](./name) to make log output easier to identify.

**Example output:**

```log
11.11 11:11	INFO	middlewares/log.go:60 abc	GET /v2/send
```
