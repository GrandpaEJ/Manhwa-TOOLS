package models

// TranslationBlock represents a single bounding box and its translated text
type TranslationBlock struct {
	// BoundingBox [ymin, xmin, ymax, xmax] relative to the image size
	BoundingBox [4]int `json:"bounding_box"`

	// Polygon represents the exact speech bubble contour [][x, y] on a 0-1000 scale
	Polygon [][2]int `json:"polygon"`

	// Type: "dialogue", "sfx", "monologue"
	Type string `json:"type"`

	// OriginalText is the native text extracted (for logging/debugging)
	OriginalText string `json:"original_text"`

	// TranslatedText is the final translated string
	TranslatedText string `json:"translated_text"`
}

// VisionResponse is the expected JSON response format from the API
type VisionResponse struct {
	Blocks []TranslationBlock `json:"blocks"`
}

