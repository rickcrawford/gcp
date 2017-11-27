package managers

import (
	"context"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"

	"github.com/rickcrawford/autocomplete/common"
	"github.com/rickcrawford/autocomplete/models"
)

// CatalogManager interacts with a catalog
type CatalogManager interface {
	Get(context.Context, string) (*models.Catalog, error)
	Save(context.Context, *models.Catalog) error
	Delete(context.Context, string) error
	List(context.Context, int, int) ([]models.Catalog, error)
}

var _ CatalogManager = (*catalogManager)(nil)

type catalogManager struct{}

func (catalogManager) Get(ctx context.Context, catalogID string) (*models.Catalog, error) {
	catalog := new(models.Catalog)
	key := catalogKey(ctx, catalogID)
	err := datastore.Get(ctx, key, catalog)
	if err == datastore.ErrNoSuchEntity {
		return nil, nil
	}
	return catalog, err
}

func (m catalogManager) Save(ctx context.Context, catalog *models.Catalog) error {
	now := time.Now()
	catalog.UpdatedAt = now
	if catalog.CreatedAt == common.NullTime {
		catalog.CreatedAt = now
	}

	if catalog.CatalogID == "" {
		catalog.CatalogID = uuid.NewV4().String()
	}

	if catalog.ApplicationKey == "" {
		catalog.ApplicationKey = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	}

	key := catalogKey(ctx, catalog.CatalogID)
	_, err := datastore.Put(ctx, key, catalog)
	return err
}

func (catalogManager) Delete(ctx context.Context, catalogID string) error {
	key := catalogKey(ctx, catalogID)
	err := datastore.Delete(ctx, key)
	if err != nil {
		return err
	}

	nctx, _ := appengine.Namespace(ctx, catalogID)
	index, err := search.Open(productIndexName)
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	keys := make([]*datastore.Key, 0)
	for t := index.List(nctx, &search.ListOptions{IDsOnly: true}); ; {
		var product models.SearchProduct
		productID, err := t.Next(&product)
		if err == search.Done {
			break
		}
		if err != nil {
			return err
		}
		ids = append(ids, productID)
		keys = append(keys, productKey(nctx, productID))
		if len(ids) == 200 {
			index.DeleteMulti(nctx, ids)
			datastore.DeleteMulti(nctx, keys)
			ids = make([]string, 0)
			keys = make([]*datastore.Key, 0)
		}
	}
	index.DeleteMulti(nctx, ids)
	datastore.DeleteMulti(nctx, keys)

	return nil
}

func (catalogManager) List(ctx context.Context, limit int, offset int) ([]models.Catalog, error) {
	query := datastore.NewQuery(catalogTypeName).Limit(limit).Offset(offset)
	catalogs := make([]models.Catalog, 0, limit)

	_, err := query.GetAll(ctx, &catalogs)
	if err != nil {
		return nil, err
	}

	return catalogs, nil
}

// NewCatalogManager returns a new catalog manager
func NewCatalogManager() CatalogManager {
	return catalogManager{}
}
