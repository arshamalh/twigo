# twigo (Golang Twitter Library)

<p align="center">
  <img src="./twigo.png" alt="twigo logo" width="300">
</p>

Twigo is a fast and easy to use twitter API library help you write twitter bots at any level.

## installation

```bash
go get github.com/arshamalh/twigo
```

# How to use

```go
package main

import "github.com/arshamalh/twigo"

twigo.NewClient(
  api_key = "your_api_key (also know as consumer key",
  api_secret = "your_api_secret (also know as consumer secret)",
  access_token = "access_token",
  access_secret = "access_secret",

  // You can use bearer token or four params above, both is not mandatory
  bearer_token = "bearer_token",

  // wait_on_rate_limit default is false,
  // TODO: maybe I should remove it and make default to true!?
  wait_on_rate_limit = false
)
```

## Contribution & Donation

Twigo is free and open to use for anyone, but you can donate or contribute if you like and this means a world for us.

donation url: <Not provided yet>
