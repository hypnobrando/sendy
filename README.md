# sendy

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/hypnobrando/sendy)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hypnobrando/sendy)](https://goreportcard.com/report/github.com/hypnobrando/sendy)

Go HTTP Client that prevents you from having to write boilerplate code setting up a native `*http.Client`, creating a request, and parsing the response.  This package uses the [builder pattern](https://medium.com/@haluan/golang-builder-design-pattern-a8b7c92969a7) for constructing requests and parsing responses.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/hypnobrando/sendy"
)

type (
    User struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }
)

func main() {
    var user User

    err := sendy.
        Get("https://myapi.com/users/1").
        SendIt().
        JSON(&user).
        Error()

    if err != nil {
        panic(err)
    }

    fmt.Println(user)
}
```

## Installation / Usage

To install `sendy`, use `go get`:
```
go get github.com/hypnobrando/sendy
```

Import the `hypnobrando/sendy` package into your code:
```go
import "github.com/hypnobrando/sendy"

func main() {
    httpClient := sendy.NewClient()
}
```

## Request hooks

`Hook` runs immediately before a request is sent (used by `DumpRequests`). `RequestHook` runs before **and** after the HTTP round-trip so callers can start and end spans, inject headers, or record status codes.

```go
type RequestHook interface {
    BeforeRequest(ctx context.Context, req *http.Request) (context.Context, error)
    AfterRequest(ctx context.Context, req *http.Request, resp *http.Response, err error)
}

sendy.AddRequestHook(myHook)          // every client, including Get/Post helpers
client.RequestHook(myHook)            // one client
request.RequestHook(myHook)           // one request
```

`AddRequestHook` is read when the request is sent, so it applies to clients created earlier.

## Staying Up to Date

To update `sendy` to the latest version, use `go get -u github.com/hypnobrando/sendy`.

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!
