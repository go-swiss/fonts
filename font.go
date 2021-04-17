package fonts

// Font details
type Font struct {
	// Category: The category of the font.
	Category string `json:"category,omitempty"`

	// Family: The name of the font.
	Family string `json:"family,omitempty"`

	// Files: The font files (with all supported scripts) for each one of
	// the available variants, as a key : value map.
	Files map[string]string `json:"files,omitempty"`

	// LastModified: The date (format "yyyy-MM-dd") the font was modified
	// for the last time.
	LastModified string `json:"lastModified,omitempty"`

	// Subsets: The scripts supported by the font.
	Subsets []string `json:"subsets,omitempty"`

	// Variants: The available variants for the font.
	Variants []string `json:"variants,omitempty"`

	// Version: The font version.
	Version string `json:"version,omitempty"`
}
