package models

import "time"

// Catalog represents a catalog
type Catalog struct {
	Base

	CatalogID      string `json:"catalogID"`
	ApplicationKey string `json:"applicationKey"`
}

// Category represents a root level category
type Category struct {
	Base

	CategoryID string `json:"categoryID"`
}

// Product represents an product in a catalog
type Product struct {
	Base

	ProductID    string     `json:"productID"`
	Description  string     `json:"description,omitempty"`
	Manufacturer string     `json:"manufacturer,omitempty"`
	Brand        string     `json:"brand,omitempty"`
	URL          string     `json:"url,omitempty"`
	ImageURL     string     `json:"imageUrl,omitempty"`
	Price        float64    `json:"price"`
	Shipping     float64    `json:"shipping"`
	Popularity   float64    `json:"popularity"`
	UPC          string     `json:"upc,omitempty"`
	MPN          string     `json:"mpn,omitempty"`
	SKU          string     `json:"sku,omitempty"`
	StreetDate   time.Time  `json":streetDate,omitempty"`
	Categories   []Category `json:"categories,omitempty"`
}
