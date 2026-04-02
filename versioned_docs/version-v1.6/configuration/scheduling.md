---
title: Scheduling
---

# Scheduling

Configure scheduling via the [`/v2/send`](/api/send-message) endpoint.

**Example:**

```yaml
settings:
  message:
    scheduling:
      maxHorizon: 10d
```

**Scheduling** can be disabled by setting `scheduling.enabled` to `false`.

The `maxHorizon` setting is used to set a boundary for **how far into the future** a message can be scheduled (default: **no limit**).
