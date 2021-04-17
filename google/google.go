package google

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-swiss/fonts"
)

//go:embed all
var all embed.FS

// Get the details of a font: family, category, available variants, e.t.c
// No http request is made. The font details are embeded in the package
// Returns ErrUnknownFont if the font family is not found
func GetFontDetails(ctx context.Context, family string) (*fonts.Font, error) {
	family = strings.ReplaceAll(family, " ", "")
	normalizedFamilyName := strings.ToLower(family)

	jsonFile, err := all.ReadFile(filepath.Join("all", normalizedFamilyName+".json"))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fonts.ErrUnknownFont
		}
		return nil, err
	}

	var font = new(fonts.Font)
	err = json.NewDecoder(bytes.NewBuffer(jsonFile)).Decode(font)
	if err != nil {
		return nil, err
	}

	return font, nil
}

// Get the []byte of a font variant
// ctx context.Context: To possibly cancel the http request.
// family: the font family e.g. "Open Sans".
// variant:  the font variant e.g. "700".
// cache: A cache implementation to prevent duplicate request if reusing fonts. Use nil to disable.
// Returns ErrUnknownFont if the font family is not found, or ErrMissingVariant if the font family does not contain that variant
func GetFontBytes(ctx context.Context, family string, variant string, cache fonts.Cache) ([]byte, error) {
	family = strings.ReplaceAll(family, " ", "")
	normalizedFamilyName := strings.ToLower(family)

	jsonFile, err := all.ReadFile(filepath.Join("all", normalizedFamilyName+".json"))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fonts.ErrUnknownFont
		}
		return nil, err
	}

	var webFont = new(fonts.Font)
	err = json.NewDecoder(bytes.NewBuffer(jsonFile)).Decode(webFont)
	if err != nil {
		return nil, err
	}

	theURL, ok := webFont.Files[variant]
	if !ok {
		return nil, fonts.ErrMissingVariant
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, theURL, nil)
	if err != nil {
		return nil, err
	}

	if cache != nil {
		cached, hit := cache.Get(theURL)
		if hit {
			return cached, nil
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if cache != nil {
		cache.Set(theURL, fileBytes)
	}
	return fileBytes, nil
}
