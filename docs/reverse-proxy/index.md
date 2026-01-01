---
sidebar_position: 1
title: Reverse Proxy
---

# Reverse Proxy

In this section we will be explaining why a **tls-enabled Reverse Proxy** is a must have.

## Why another Proxy

**Secured Signal API** itself is already a **Reverse Proxy**, lacking one important feature: **SSL Certificates**.

### SSL Certificates

If you want to deploy anything on the Internet a **SSL Certificate** is almost a necessity, same goes for **Secured Signal API**,
even if you don't plan on exposing your instance to the Internet it is always good to have an extra layer of **Security**,
**SSL Certificates** are needed for establishing **secure HTTP requests**.

### Port forwarding

Furthermore if you want to have multiple services **on the same port** using **HTTP** you'd also need a **tls-enabled Reverse Proxy**,
to route requests to the correct backend based on hostnames and routing rules.

### Not Convinced?

And if you are still not convinced then look at this [article](https://www.cloudflare.com/learning/cdn/glossary/reverse-proxy) online.
