package models

import (
	"time"

	"google.golang.org/appengine/search"

	"github.com/rickcrawford/autocomplete/common"
)

// SearchQuery is a query to perform
type SearchQuery struct {
	Limit int
	Query string
}

// SearchProduct implements the search loader interface
type SearchProduct struct {
	Product
}

// Load is for search indexing
func (p *SearchProduct) Load(fields []search.Field, meta *search.DocumentMetadata) error {
	// Load the title field, failing if any other field is found.
	for _, f := range fields {
		switch f.Name {
		case "ProductID":
			p.ProductID = f.Value.(string)
		case "Name":
			p.Name = f.Value.(string)
		case "Manufacturer":
			p.Manufacturer = f.Value.(string)
		case "Brand":
			p.Brand = f.Value.(string)
		case "Popularity":
			p.Popularity = f.Value.(float64)
		case "Price":
			p.Price = f.Value.(float64)
		case "CreatedAt":
			p.CreatedAt = f.Value.(time.Time)
		case "UpdatedAt":
			p.UpdatedAt = f.Value.(time.Time)
		}
	}
	return nil
}

// Save is for search indexing.
func (p *SearchProduct) Save() ([]search.Field, *search.DocumentMetadata, error) {
	fields := []search.Field{
		{Name: "ProductID", Value: p.ProductID},
		{Name: "Name", Value: p.Name},
		{Name: "Popularity", Value: p.Popularity},
		{Name: "Price", Value: p.Price},
		{Name: "CreatedAt", Value: p.CreatedAt},
		{Name: "UpdatedAt", Value: p.UpdatedAt},
	}

	prefix := common.FormatPrefix(p.Name, "_")
	for i := range prefix {
		fields = append(fields, search.Field{Name: "prefix", Value: prefix[0:i]})
		if i == common.MaxPrefixLength {
			break
		}
	}
	fields = append(fields, search.Field{Name: "prefix_sort", Value: prefix})

	return fields, nil, nil
}
