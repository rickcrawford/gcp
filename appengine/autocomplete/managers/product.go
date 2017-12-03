package managers

import (
	"context"
	"errors"
	"strconv"
	"time"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"

	"github.com/rickcrawford/gcp/appengine/autocomplete/common"
	"github.com/rickcrawford/gcp/common/models"
)

const MaxPrefixLength = 10

// SearchQuery is a query to perform
type SearchQuery struct {
	Limit int
	Query string
}

// SearchProduct implements the search loader interface
type SearchProduct struct {
	SKU  string
	Name string
}

// Load is for search indexing
func (p *SearchProduct) Load(fields []search.Field, meta *search.DocumentMetadata) error {
	// Load the title field, failing if any other field is found.
	for _, f := range fields {
		switch f.Name {
		case "SKU":
			p.SKU = f.Value.(string)
		case "Name":
			p.Name = f.Value.(string)
		}
	}
	return nil
}

// Save is for search indexing.
func (p *SearchProduct) Save() ([]search.Field, *search.DocumentMetadata, error) {
	fields := []search.Field{
		{Name: "SKU", Value: p.SKU},
		{Name: "Name", Value: p.Name},
	}

	prefix := common.FormatPrefix(p.Name, "_")
	for i := range prefix {
		fields = append(fields, search.Field{Name: "prefix", Value: prefix[0:i]})
		if i == MaxPrefixLength {
			break
		}
	}
	fields = append(fields, search.Field{Name: "prefix_sort", Value: prefix})

	return fields, nil, nil
}

const maxProducts = 200

// ErrTooManyProducts returned when there are too many products in the response
var ErrTooManyProducts = errors.New("too many products")

// ProductManager gets product data
type ProductManager interface {
	Get(context.Context, int) (*models.Product, error)
	Delete(context.Context, int) error
	Save(context.Context, *models.Product) error
	SaveAll(context.Context, []models.Product) error
	List(context.Context, int, int) ([]models.Product, error)
	Search(context.Context, SearchQuery) ([]models.Product, error)
}

var _ ProductManager = (*productManager)(nil)

type productManager struct{}

func (productManager) Get(ctx context.Context, SKU int) (*models.Product, error) {
	product := new(models.Product)
	key := productKey(ctx, SKU)
	err := datastore.Get(ctx, key, product)
	if err == datastore.ErrNoSuchEntity {
		return nil, nil
	}
	return product, err
}

func (productManager) Delete(ctx context.Context, SKU int) error {
	key := productKey(ctx, SKU)
	return datastore.Delete(ctx, key)
}

func (productManager) SaveAll(ctx context.Context, products []models.Product) error {
	if len(products) > maxProducts {
		return ErrTooManyProducts
	}

	index, err := search.Open(productIndexName)
	if err != nil {
		return err
	}

	now := time.Now()
	keys := make([]*datastore.Key, len(products))
	ids := make([]string, len(products))
	searchProducts := make([]interface{}, len(products))

	for i, product := range products {
		keys[i] = productKey(ctx, product.SKU)
		ids[i] = strconv.Itoa(product.SKU)

		products[i].Updated = now
		searchProducts[i] = &SearchProduct{
			SKU:  strconv.Itoa(products[i].SKU),
			Name: products[i].Name,
		}
	}
	if _, err = datastore.PutMulti(ctx, keys, products); err != nil {
		return err
	}

	_, err = index.PutMulti(ctx, ids, searchProducts)
	return err
}

func (p productManager) Save(ctx context.Context, product *models.Product) error {
	return p.SaveAll(ctx, []models.Product{*product})
}

func (productManager) List(ctx context.Context, limit int, offset int) ([]models.Product, error) {
	query := datastore.NewQuery(productTypeName).Limit(limit).Offset(offset)
	products := make([]models.Product, 0, limit)

	_, err := query.GetAll(ctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (productManager) Search(ctx context.Context, query SearchQuery) ([]models.Product, error) {
	index, err := search.Open(productIndexName)
	if err != nil {
		return nil, err
	}

	products := make([]models.Product, 0)
	result := index.Search(ctx, query.Query, &search.SearchOptions{
		Limit: query.Limit,
		Sort: &search.SortOptions{
			Expressions: []search.SortExpression{{Expr: "prefix_sort", Reverse: true}},
		},
	})

	for {
		var product SearchProduct
		_, err = result.Next(&product)
		if err == search.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		SKU, _ := strconv.Atoi(product.SKU)
		products = append(products, models.Product{SKU: SKU, Name: product.Name})
	}

	return products, nil
}

// NewProductManager gets a new product manager
func NewProductManager() ProductManager {
	return productManager{}
}
