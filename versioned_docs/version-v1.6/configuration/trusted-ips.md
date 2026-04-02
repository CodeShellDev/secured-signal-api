---
title: Trusted IPs
---

# Trusted IPs

Trusted clients can bypass some security features and are often local or internal IPs.

To trust IPs or ranges add them to `access.trustedIPs`.

```yaml
settings:
  access:
    trustedIPs:
      - 192.168.1.10
      - 192.168.2.0/24
```
