---
sidebar_position: 5
title: Best Practices
---

# Best Practices

Here are some common best practices for running **Secured Signal API**, but these generally apply for any service.

## Usage

- Create **separate configs** for each service
- Use [**placeholders**](./usage/advanced#placeholders) extensively
- Always keep your stack **up-to-date**

## Security

- Always use **API tokens** in production
- Run behind a [**tls-enabled reverse proxy**](./reverse-proxy)
- Be cautious when overriding [**blocked endpoints**](./features#endpoints)
- Use per-token overrides to **enforce least privilege**
- Always allow the least possible [**endpoints**](./features#endpoints)
- Only allow access from [**IPs**](./features#ip-filter) **you trust**
