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

## Staying Up to Date

To update `sendy` to the latest version, use `go get -u github.com/hypnobrando/sendy`.

## Contributing

Please feel free to submit issues, fork the repository and send pull requests!
