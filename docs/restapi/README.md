# REST API Docs

When enabled by config, shinpuru exposes a RESTful HTTP API which exposes all functionalities which are also available to the web frontend.

## Authentication

All requests to the API needs to be authenticated and authorized. To authenticate your requests, you need to generate an API token in shinpurus web interface.

![](https://i.imgur.com/KYp2OdR.png)
![](https://i.imgur.com/RBrQwrH.png)
![](https://i.imgur.com/XPS0h7R.png)

To authenticate your requests, you need to add an `Authentication` header to your request with the token as `Bearer` value.

```
> GET /api/me HTTP/2.0
> Host: shnp.de
> Authorization: bearer eyJhbGciOiJIUzI1...
> Accept: */*
```

## Endpoints

The endpoints for the V1 API can be found [**in this document**](v1/restapi.md).

Alternatively, you can find [**here**](https://app.swaggerhub.com/apis-docs/zekroTJA/shinpuru-main-api/1.0) a more interactive representation of the REST API documentation.