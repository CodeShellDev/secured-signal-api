---
title: Log Levels
---

# Log Levels

Log Levels are often used in programms to **filter out to verbose** information or to allow for debugging logs.

To change the Log Level set `logLevel` to: (default: `info`)

**Levels:**

- `info`
- `debug` (verbose)
- `warn` (**only** warnings and errors)
- `error` (**only** errors)
- `fatal` (**only** fata errors)

> [!CAUTION]
> the log level `dev` **can leak data in the logs**
> and must only be used for testing during development
