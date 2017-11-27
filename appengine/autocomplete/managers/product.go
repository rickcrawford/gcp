package managers

import (
	"context"
	"errors"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"

	"github.com/rickcrawford/autocomplete/common"
	"github.com/rickcrawford/autocomplete/models"
)

const maxProducts = 200

// ErrTooManyProducts returned when there are too many products in the response
var ErrTooManyProducts = errors.New("too many products")

// ProductManager gets product data
type ProductManager interface {
	Get(context.Context, string, string) (*models.Product, error)
	Delete(context.Context, string, string) error
	Save(context.Context, string, *models.Product) error
	SaveAll(context.Context, string, []models.Product) error
	List(context.Context, string, int, int) ([]models.Product, error)
	Search(context.Context, string, models.SearchQuery) ([]models.Product, error)
}

var _ ProductManager = (*productManager)(nil)

type productManager struct{}

func (productManager) Get(ctx context.Context, catalogID, productID string) (*models.Product, error) {
	nctx, _ := appengine.Namespace(ctx, catalogID)

	product := new(models.Product)
	key := productKey(nctx, productID)
	err := datastore.Get(nctx, key, product)
	if err == datastore.ErrNoSuchEntity {
		return nil, nil
	}
	return product, err
}

func (productManager) Delete(ctx context.Context, catalogID, productID string) error {
	nctx, _ := appengine.Namespace(ctx, catalogID)

	key := productKey(nctx, productID)
	return datastore.Delete(nctx, key)
}

func (productManager) SaveAll(ctx context.Context, catalogID string, products []models.Product) error {
	if len(products) > maxProducts {
		return ErrTooManyProducts
	}

	index, err := search.Open(productIndexName)
	if err != nil {
		return err
	}

	nctx, _ := appengine.Namespace(ctx, catalogID)

	now := time.Now()
	keys := make([]*datastore.Key, len(products))
	ids := make([]string, len(products))
	searchProducts := make([]interface{}, len(products))

	for i, product := range products {
		keys[i] = productKey(nctx, product.ProductID)
		ids[i] = product.ProductID

		products[i].UpdatedAt = now
		if products[i].CreatedAt == common.NullTime {
			products[i].CreatedAt = now
		}
		searchProducts[i] = &models.SearchProduct{products[i]}
	}
	if _, err = datastore.PutMulti(nctx, keys, products); err != nil {
		return err
	}

	_, err = index.PutMulti(nctx, ids, searchProducts)
	return err
}

func (p productManager) Save(ctx context.Context, catalogID string, product *models.Product) error {
	return p.SaveAll(ctx, catalogID, []models.Product{*product})
}

func (productManager) List(ctx context.Context, catalogID string, limit int, offset int) ([]models.Product, error) {
	nctx, _ := appengine.Namespace(ctx, catalogID)

	query := datastore.NewQuery(productTypeName).Limit(limit).Offset(offset)
	products := make([]models.Product, 0, limit)

	_, err := query.GetAll(nctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (productManager) Search(ctx context.Context, catalogID string, query models.SearchQuery) ([]models.Product, error) {
	index, err := search.Open(productIndexName)
	if err != nil {
		return nil, err
	}

	nctx, _ := appengine.Namespace(ctx, catalogID)

	products := make([]models.Product, 0)
	result := index.Search(nctx, query.Query, &search.SearchOptions{
		Limit: query.Limit,
		Sort: &search.SortOptions{
			Expressions: []search.SortExpression{{Expr: "prefix_sort", Reverse: true}},
		},
	})

	for {
		var product models.SearchProduct
		_, err = result.Next(&product)
		if err == search.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		products = append(products, product.Product)
	}

	return products, nil
}

// NewProductManager gets a new product manager
func NewProductManager() ProductManager {
	return productManager{}
}
