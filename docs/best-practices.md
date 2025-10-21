---
sidebar_position: 4
title: Best Practices
---

# Best Practices

Here are some common best practices for running **Secured Signal API**, but these generally apply for any service.

## Usage

- Create **seperate configs** for each service
- Use **Placeholders** extensively _(they are your friends)_
- Always keep your stack **up-to-date** _(this is why we have docker)_

## Security

- Always use **API tokens** in production
- Run behind a **tls-enabled** [Reverse Proxy](./reverse-proxy/overview)
- Be cautious when overriding **Blocked Endpoints**
- Use per-token overrides to **enforce least privilege**
- Always allow the least possible access points
