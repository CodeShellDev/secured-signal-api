---
sidebar_position: 1
title: Reverse Proxy
---

# Reverse Proxy

In this section we will be explaining why a **tls-enabled reverse proxy** is a must-have.

## Why another Proxy

**Secured Signal API** itself is already a **reverse proxy**, lacking one important feature: **SSL certificates**.

### SSL Certificates

When deploying anything on the internet an **SSL certificate** is almost a necessity.
Same goes for **Secured Signal API**, even if you don't plan on exposing your instance to the internet it is always good to have an extra layer of **security**.

### Port Forwarding

Furthermore, if you want to have multiple services **on the same port** using **HTTP** you'd also need a **tls-enabled reverse proxy**,
to route requests to the correct backend based on hostnames and routing rules.

### Not Convinced?

And if you are still not convinced then look at this [article](https://www.cloudflare.com/learning/cdn/glossary/reverse-proxy) online.
