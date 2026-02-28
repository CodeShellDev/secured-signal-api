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
      "@number":
        - value: "+123400002"
          action: block
        - value: "+12340000[1-9]"
          matchType: regex
          action: allow
```

Set the wanted action on encounter, available options are `block` and `allow`.

> [!NOTE]
> Supported [placeholder types](../usage/advanced#placeholders):
>
> | `.` Variables | `@` Body | `#` Headers |
> | ------------- | -------- | ----------- |
> | ❌            | ✅       | ✅          |

## Match Types

Available options for `matchType` are:

{{{ #://../templates/match-rules.md.tmpl }}}

## Behavior

| Allow               | Block                             | Result                                                                      |
| ------------------- | --------------------------------- | --------------------------------------------------------------------------- |
| `number=+123400003` | —                                 | `number` may **only** be `+123400003`                                       |
| —                   | `number=+123400002`               | `number` may **not** be `+123400002`                                        |
| `message=hello`     | `number=+123400002`               | `number` may **not** be `+123400002`<br/> `message` may **only** be `hello` |
| `number=+123400003` | `number=+12340000[1-9]` (`regex`) | `number` may **not** be `+123400001` through `9` **except** `123400003`     |

### Rules

- **Field-scoped** (policies for `a` don't affect policies for `b`)

{{{ #://../templates/block-rules.md.tmpl }}}
