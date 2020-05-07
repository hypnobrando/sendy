# sendy

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/brandoneprice31/sendy)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/brandoneprice31/sendy)](https://goreportcard.com/report/github.com/brandoneprice31/sendy)

## Quick Start

```go
package main

import (
    "github.com/brandoneprice31/sendy"
)

type (
    User struct {
        ID   int    `json:"id""`
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
