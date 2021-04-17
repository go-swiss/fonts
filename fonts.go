package fonts

import "errors"

//go:generate go run ./generate

var (
	ErrUnknownFont    error = errors.New("Unknown font")
	ErrMissingVariant error = errors.New("Missing variant")
)

// Cache can be passed to certain functions to prevent duplicate requests if reusing fonts
type Cache interface {
	Get(string) (value []byte, hit bool)
	Set(key string, value []byte)
}
