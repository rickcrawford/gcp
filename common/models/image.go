package models

type Image struct {
	SKU     int    `json:"sku"`
	Image   string `json:"image"`
	Content string `json:"content"`
}
