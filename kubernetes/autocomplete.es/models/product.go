package models

import "time"

// Category represents a category
type Category struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// Suggestion field is used for generating suggestions
type Suggestion struct {
	Input  string `json:"input"`
	Weight int    `json:"weight"`
}

// Product represents a product
type Product struct {
	SKU          int        `json:"sku"`
	Name         string     `json:"name"`
	Price        float64    `json:"price"`
	UPC          string     `json:"upc"`
	Category     []Category `json:"category"`
	Description  string     `json:"description"`
	Manufacturer string     `json:"manufacturer"`
	Model        string     `json:"model"`
	URL          string     `json:"url"`
	Image        string     `json:"image"`
	Content      string     `json:"content"`
	Updated      time.Time  `json:"updated"`

	// the following fields are not part of the feed but can be used to influence results
	Suggestion []Suggestion `json:"suggestion,omitempty"`
	Weight     float64      `json:"weight"`
}
