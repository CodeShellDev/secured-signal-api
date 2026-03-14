+{{{ import "~/docs/templates/functions.inc.gtmpl" }}}
<img align="center" width="1048" height="512" alt="Secure API Gateway Proxy for Signal CLI REST API" src="https://github.com/codeshelldev/secured-signal-api/raw/refs/heads/docs/static/img/banner.png" />

<h3 align="center">Secure API Gateway Proxy for <a href="https://github.com/bbernhard/signal-cli-rest-api">Signal CLI REST API</a></h3>

<p align="center">
token-based authentication,
endpoint restrictions, placeholders, flexible configuration
</p>

<p align="center">
🔒 Secure · ⭐️ Configurable · 🚀 Easy to Deploy with Docker
</p>

<div align="center">
  <a href="https://github.com/codeshelldev/secured-signal-api/releases">
    <img 
		src="https://img.shields.io/github/v/release/codeshelldev/secured-signal-api?sort=semver&logo=github&label=Release" 
		alt="GitHub release">
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/stargazers">
    <img 
		src="https://img.shields.io/github/stars/codeshelldev/secured-signal-api?style=flat&logo=github&label=Stars" 
		alt="GitHub stars">
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/pkgs/container/secured-signal-api">
    <img
    src='https://img.shields.io/badge/Image%20Size-+{{{ replace ( htmlText ( htmlDocFind ( htmlDecode ( fetch "https://ghcr-badge.egpl.dev/codeshelldev/secured-signal-api/size?color=%2344cc11&tag=latest&label=Image+Size&trim=" ) ) "svg g:nth-of-type(3) text" ) ) " " "%20" }}}-_?color=2344cc11'
    alt="Docker image Pulls">
  </a>
  <a href="https://github.com/codeshelldev/secured-signal-api/pkgs/container/secured-signal-api">
    <img 
		src="https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fghcr-badge.elias.eu.org%2Fapi%2Fcodeshelldev%2Fsecured-signal-api%2Fsecured-signal-api&query=downloadCount&label=Pulls&color=2344cc11&logo=docker"
		alt="Docker image Pulls">
  </a>
  <a href="./LICENSE">
    <img 
		src="https://img.shields.io/badge/License-MIT-green.svg"
		alt="License: MIT">
  </a>
</div>

## Contents

> [!IMPORTANT]
> Check out the [**Official Documentation**](https://codeshelldev.github.io/secured-signal-api) for up-to-date instructions and additional content!

> [!WARNING]
> We are slowly moving away from this README and instead are trying to make the [**Official Documentation**](https://codeshelldev.github.io/secured-signal-api) the only source of truth

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Features](#features)
- [Configuration](#configuration)
  - [Endpoints Restrictions](https://codeshelldev.github.io/secured-signal-api/docs/configuration/endpoints)
  - [Placeholder Variables](https://codeshelldev.github.io/secured-signal-api/docs/configuration/variables)
  - [Field Policies](https://codeshelldev.github.io/secured-signal-api/docs/configuration/field-policies)
  - [Field Mappings](https://codeshelldev.github.io/secured-signal-api/docs/configuration/field-mappings)
  - [Message Template](https://codeshelldev.github.io/secured-signal-api/docs/configuration/message-template)
  - [Port Restrictions](https://codeshelldev.github.io/secured-signal-api/docs/configuration/port)
  - [Hostname Restrictions](https://codeshelldev.github.io/secured-signal-api/docs/configuration/hostnames)
  - [IP Filter](https://codeshelldev.github.io/secured-signal-api/docs/configuration/ip-filter)
  - [Rate Limiting](https://codeshelldev.github.io/secured-signal-api/docs/configuration/rate-limiting)
  - [Trusted Proxies](https://codeshelldev.github.io/secured-signal-api/docs/configuration/trusted-proxies)
  - [Trusted IPs](https://codeshelldev.github.io/secured-signal-api/docs/configuration/trusted-ips)
  - [Log Level](https://codeshelldev.github.io/secured-signal-api/docs/configuration/log-level)
  - [Port](https://codeshelldev.github.io/secured-signal-api/docs/configuration/port)
  - [Hostnames](https://codeshelldev.github.io/secured-signal-api/docs/configuration/hostnames)
  - [Auth Methods](https://codeshelldev.github.io/secured-signal-api/docs/configuration/auth)
- [Reverse Proxy](https://codeshelldev.github.io/secured-signal-api/docs/reverse-proxy)
- [Integrations](https://codeshelldev.github.io/secured-signal-api/docs/integrations)
- [Contributing](#contributing)
- [Support](#support)
- [Help](#help)
- [License](#license)
- [Legal](#legal)

## Getting Started

+{{{ funcCallArgs "parseDocs" "~/docs/getting-started/installation.md" "### Installation" }}}

### Setup

Once you have installed **Secured Signal API** you can [register or link a signal account](https://codeshelldev.github.io/secured-signal-api/docs/getting-started/setup).

+{{{ funcCallArgs "parseDocs" "~/docs/usage/index.md" "## Usage" }}}

+{{{ funcCallArgs "parseDocs" "~/docs/features/features.md" "## Features" }}}

+{{{ funcCallArgs "parseDocs" "~/docs/configuration/index.md" "## Configuration" }}}

## Contributing

Found a bug? Want to change or add something?
Feel free to open up an [issue](https://github.com/codeshelldev/secured-signal-api/issues) or create a [pull request](https://github.com/codeshelldev/secured-signal-api/pulls)!

## Support

Has this Repo been helpful 👍️ to you? Then consider ⭐️'ing this Project.

:)

## Help

**Are you having problems setting up Secured Signal API?**<br/>
No worries check out the [discussions](https://github.com/codeshelldev/secured-signal-api/discussions) tab and ask for help.

**We are all volunteers**, so please be friendly and patient.

## License

This Project is licensed under the [MIT License](./LICENSE).

## Legal

Logo designed by [@CodeShellDev](https://github.com/codeshelldev), All Rights Reserved.

This Project is not affiliated with the Signal Foundation.
