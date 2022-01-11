# REST API Docs

When enabled by config, shinpuru exposes a RESTful HTTP API which exposes all functionalities which are also available to the web frontend.

## Authentication

All requests to the API needs to be authenticated and authorized. To authenticate your requests, you need to generate an API token in shinpuru's web interface.

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

## Endpoints and Models

The endpoint and model documentation for the V1 API can be found [**in this document**](v1/restapi.md).

Alternatively, you can find [**here**](https://app.swaggerhub.com/apis-docs/zekroTJA/shinpuru-main-api/1.0) a more interactive representation of the REST API documentation.

## Some Things to know before using the API

### Data Update Behavior

Because Go has [default zero values](https://tour.golang.org/basics/12) for each primitive type like integers, strings and bools, you can not determine that easily if a value is meant to be actually `0`, `false` or an empty string (`""`), because these are also the default values when nothing is specified.

There are two ways around this.

1. Use pointers for everything.
2. Specify that **all** values are defined as "set".

The first solution was not suitable in my opinion, because it would reqire a lot code around `nil` and proper value checking of each model property, which would also introduce a lot new fault sources. Also, because shinpuru's API utilizes a lot of the original models of [discordgo](https://github.com/bwmarrin/discordgo), this would require a lot of model wrapping and double definitions.

So, I went for the second solution*. 

Every property is specified as "set" by the API on update. That means, if you pass `null` as value of a string, that means the valaue of the property will be updated to `""`, which is the default zero value of a string. Long story short, even if you want to update only single properties of a model, you must pass the whole model on update to ensure consistency, even if this means that you need to get the current values before you can update them.

Let's take the [`/settings/presence` endpoint](), for example. We want to update the `game` value.

```javascript
// Perform a request with the given method to the passed
// uri with the given data (on)
async function request(method, uri, data) {
  const res = await window.fetch(uri, { 
    method,
    body: data ? JSON.stringify(data) : null,
    headers: { /* Auth stuff ... */ },
  });
  return await res.json();
}

(async () => {
  // Request the current presence state.
  const presence = await request('GET', 'https://shnp.de/api/v1/settings/presence');
  // Update the game status of the presence.
  presence.game = 'https://c.shnp.de';
  // Update the presence status with the whole presence state object.
  await request('POST', 'https://shnp.de/api/v1/settings/presence', presence);
})();
```

