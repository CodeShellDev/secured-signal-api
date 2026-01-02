---
title: Field Policies
---

# Field Policies

An extra layer of security for ensuring no unwanted values are passed through a request.

**Field Policies** allow for blocking or specifically allowing certain fields with set values from being used in the requests body or headers.

Configure them by using `access.fieldPolicies` like so:

```yaml
settings:
  access:
    fieldPolicies:
      "@number": { value: "+123400002", action: block }
```

Set the wanted action on encounter, available options are `block` and `allow`.

> [!IMPORTANT]
> Supported [placeholder types](../usage/advanced#placeholders):
>
> | `.` Variables | `@` Body | `#` Headers |
> | ------------- | -------- | ----------- |
> | ❌            | ✅       | ✅          |
