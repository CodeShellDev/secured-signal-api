---
title: Log Level
---

# Log Level

Log levels help to filter or explicitly allow certain information from the logs.

To change the log level set `logLevel` to: (default: `info`)

**Levels:**

- `info`
- `debug` (verbose)
- `warn` (**only** warnings and errors)
- `error` (**only** errors)
- `fatal` (**only** fatal errors)

> [!CAUTION]
> The log level `dev` **can leak data in the logs**
> and must only be used for testing during development
