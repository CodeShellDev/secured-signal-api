---
title: Trusted Proxies
---

# Trusted Proxies

Proxies can be marked as trusted.

When determining the IP of a client behind a proxy it is important to use the `X-Forwarded-For` header,
but as you might know anyone can set headers (spoofing possible).

To prevent IP spoofing you should only trust the HTTP headers of trusted proxies.
Otherwise malicious actors may change the `X-Forwarded-For` header to be able to bypass block list or rate limits.

Add proxies to be trusted in `access.trustedProxies`.
