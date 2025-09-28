#### NGINX Proxy

This is the [NGINX](https://github.com/nginx/nginx) `docker-compose.yaml` file:

```yaml
{ { file.examples/reverse-proxy/nginx/nginx.docker-compose.yaml } }
```

Create a `nginx.conf` file in the `docker-compose.yaml` folder and mount it to `etc/nginx/conf.d/default.conf`:

```conf
{ { file.examples/reverse-proxy/nginx/nginx.conf } }
```

Lastly add your `cert.key` and `cert.crt` into your `certs/` folder and mount it to `/etc/nginx/ssl`.
