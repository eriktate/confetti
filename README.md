# confetti

Automagically load configuration directly into structs.

## Installation

```bash
go get github.com/eriktate/confetti
```

## Getting Started

`confetti` uses reflection to map your struct fields to environment variable keys and
values. By default the field names are used to match keys, but in practice you'll want
to use the `conf` field attribute to re-map them. 

Here's an example of how `confetti` works:
```go
package main

import (
    "fmt"

    "github.com/eriktate/confetti"
)

type Config struct {
    SomeString string `conf:"MY_STRING"`
    SomeInt int `conf:"MY_INT"`
}

func main() {
    cfg, err := confetti.FromEnv[Config]()
    if err != nil {
        fmt.Fatalf("failed to load config: %s", err)
    }

    fmt.Printf("SomeString=%s SomeInt=%d", cfg.SomeString, cfg.SomeInt)
}
```

If you run the example with valid env vars for `MY_STRING` and `MY_INT`, you should see
both values printed. You can also load configuration from `.env` formatted files and
overaly configurations. For example, if you wanted to load a `.env` file, a `.secret`
file, and allow for overrides through environment variables you could easily do so:

```go
package main

import (
    "fmt"

    "github.com/eriktate/confetti"
)

type Config struct {
    ClientID string `conf:"MY_CLIENT_ID"`
    ClientSecret int `conf:"MY_CLIENT_SECRET"`
}

func main() {
    cfg, err := confetti.FromFiles[Config](".env", ".secret")
    if err != nil {
        fmt.Fatalf("failed to load config from files: %s", err)
    }

    if err := ApplyEnv(cfg); err != nil {
      fmt.Fatalf("failed to load config from env vars: %s", err)
    }

    fmt.Printf("ClientID=%s ClientSecret=%d", cfg.ClientID, cfg.ClientSecret)
}
```
This will hydrate your `Config` struct using the configuration found in `.env`, and
`.secret` while preferring values in `.secret` if there are any overlaps. Values provided
as environment variables will take ultimate precedence since they're applied to the
`Config` struct last.

## Why build this?

I don't like pulling in random dependencies for simple things I could write for myself in
an evening, so I put together `confetti` in order to have a flexible way of loading
configuration for my side projects. Since this repo becomes a usable go library anyway, I
figured I'd at least throw in a README in case some poor soul stumbles upon it
and wants to use it. Yes, I'm fully aware of the irony.

## Is this maintained?

If there's something that makes sense to add, I'll add it. If there are bugs reported,
I'll fix them. There's a much more fully-featured version of `confetti` that I could
build, but for now I'm just adding the features I personally need.

