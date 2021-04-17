# Fonts

Fonts is a package that provides helpers to access font details and easily retrieve font bytes.

This package has **ZERO** 3rd-party dependencies.

For now, Google Fonts is supported, but other fonts are planned.

## Motivation

When working with images, PDFs or any other format that requires us to embed fonts, we can decide to download the fonts and embed them directly into our application.

However, if our application could use multiple different fonts which are decided at runtime, then there are a couple complications:

1. There is no easy way to know what fonts are available and what variants they have. Google Fonts provides an API, but you can only receive a dump all **FONT** details.
2. It becomes tedious to download the font bytes during runtime, and manage caching.

This package helps to simplify this process

## Reference

See <https://pkg.go.dev/github.com/go-swiss/fonts>

## How to use

### Get Font Details

To get the details of a Google Font, use the `GetFontDetails` function:

```go
package main 

import (
    "fmt"
    "context"

    "github.com/go-swiss/fonts/google"
)

func main() {
    ctx := context.Background()
    details, err := google.GetFontDetails(ctx, "Roboto")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Family: %s\nVariants: %v", details.Family, details.Variants)
}
```

```shell
Family: Roboto
Variants: [100 100italic 300 300italic regular italic 500 500italic 700 700italic 900 900italic]
```

### Get Font Bytes

To get the `bytes` of a font, use the `GetFontBytes` method. This will download the font bytes from the appropriate URL.

```go
family := "Open Sans"
variant := "900"

fontBytes, err := google.GetFontBytes(ctx, family, variant, nil)
if errors.Is(err, fonts.ErrMissingVariant) {
    fontBytes, err = google.GetFontBytes(ctx, family, "regular", nil)
}

if err != nil {
    panic(err)
}

font, err := opentype.Parse(fontBytes)
if err != nil {
    panic(err)
}

// Do something with the font
```

### Caching Responses

In our application it is possible that we need to frequetnly get fonts. To reduce duplicate requests, a `fonts.Cache` implementation can be passed as the 4th parameter to the `GetFontBytes` function.

Here is an example using <https://github.com/ReneKroon/ttlcache>


```go
package main 

import (
    "time"

    "github.com/ReneKroon/ttlcache/v2"
)

type cache struct {
	c *ttlcache.Cache
}

func (c cache) Get(key string) ([]byte, bool) {
	fontInterface, err := c.c.Get(key)
	if err != nil {
		return nil, false
	}

	fontBytes, ok := fontInterface.([]byte)
	return fontBytes, ok
}

func (c cache) Set(key string, val []byte) {
	c.c.Set(key, val)
}

func main() {
	c := ttlcache.NewCache()
	c.SetTTL(time.Duration(20 * time.Second))
	c.SetCacheSizeLimit(64)

    fontBytes, err := google.GetFontBytes(config.Context, family, variant, cache{c})
    if err != nil {
        panic(err)
    }

    font, err := opentype.Parse(fontBytes)
    if err != nil {
        panic(err)
    }

    // Do something with the font
}
```

## Contributing

Feature requests and Pull Requests are welcome!

If working locally, the `google/all` directory can be refreshed by running `go generate`. See `generate/main.go` for details.
