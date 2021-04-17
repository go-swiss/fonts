package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-swiss/fonts"
)

const googleWebfontsListURL = "https://www.googleapis.com/webfonts/v1/webfonts"

type fontList struct {
	// Items: The list of fonts currently served by the Google Fonts API.
	Items []*fonts.Font `json:"items,omitempty"`

	// Kind: This kind represents a list of webfont objects in the webfonts
	// service.
	Kind string `json:"kind,omitempty"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	apikey := os.Getenv("GOOGLE_FONTS_API_KEY")
	if apikey == "" {
		log.Fatalf("env var GOOGLE_FONTS_API_KEY is needed to generate files")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleWebfontsListURL+"?key="+apikey, nil)
	if err != nil {
		log.Fatalf("could not create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("could not do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad Status Code: %#v", resp)
	}

	var fonts fontList
	err = json.NewDecoder(resp.Body).Decode(&fonts)
	if err != nil {
		log.Fatalf("could not decode fonts response: %v", err)
	}

	err = generate(ctx, fonts)
	if err != nil {
		log.Fatalf("could not generate constants: %v", err)
	}
}

func generate(ctx context.Context, fonts fontList) error {
	for _, font := range fonts.Items {
		err := addFontJSON(font)
		if err != nil {
			return fmt.Errorf("could not add JSON file: %w", err)
		}
	}

	return nil
}

func addFontJSON(font *fonts.Font) error {
	normalizedFamilyName := strings.ToLower(strings.ReplaceAll(font.Family, " ", ""))

	log.Printf("Adding google/all/%s.json", normalizedFamilyName)
	jsonFileName := filepath.Join("google", "all", normalizedFamilyName+".json")
	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		return fmt.Errorf("could not create file %q: %w", jsonFileName, err)
	}
	defer jsonFile.Close()

	err = json.NewEncoder(jsonFile).Encode(font)
	if err != nil {
		return fmt.Errorf("could encode font JSON: %w", err)
	}

	return nil
}
