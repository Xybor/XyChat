# DEPLOY A SECURE WEB APPLICATION
## Step 1. Create a certificate
In this guide, we help you create a self-signed certificate using openssl.  You
could get a trusted certificate by buying it at hosting companies.

OpenSSL is available on either Windows or Linux OS. For windows, download the OpenSSL [here](https://slproweb.com/products/Win32OpenSSL.html).

```shell
Xychat$ openssl genrsa -out xychat.key 2048
Xychat$ openssl ecparam -genkey -name secp384r1 -out xychat.key
Xychat$ openssl req -new -x509 -sha256 -key xychat.key -out xychat.crt -days 3650
```

## Step 2. Set the environtment variable `TLS`:
```shell
Xychat$ set TLS=xychat
```

## Note
If you don't want to continue to use HTTPS protocol, let you unset the `TLS` environment variable.
```shell
Xychat$ set TLS=
```
