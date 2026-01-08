---
title: Port
---

# Port

To change the port which **Secured Signal API** uses, you need to set `service.port` in your config. (default: `8880`)

## Token-specific ports (port realms)

You can additionally define a port per Token config.

> [!NOTE]
> Each port spawns a separate listener, beware that this _can_ mildly affect performance

When a token specifies a port, a new **realm** is created for that port.  
Only tokens that explicitly belong to the same realm are accepted on that port.

### Example

- `TOKEN_1` → port `8880`
- `TOKEN_2` → port `8890`

Requests behave as follows:

| Token     | Port   | Result |
| :-------- | :----- | :----: |
| `TOKEN_1` | `8880` |   ✅   |
| `TOKEN_2` | `8880` |  ⛔️   |
| `TOKEN_1` | `8890` |  ⛔️   |
| `TOKEN_2` | `8890` |   ✅   |

If a token config does not specify a port, it automatically gets assigned to the default realm.

This allows strict separation of access by port without running multiple instances.
