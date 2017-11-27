package models

import "time"

// Base is a base object
type Base struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Disabled  bool      `json:"disabled"`

	Batch string `json:"batch,omitempty"`
}
