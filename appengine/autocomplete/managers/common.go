package managers

import (
	"context"

	"google.golang.org/appengine/datastore"
)

const catalogTypeName = "Catalog"
const productTypeName = "Product"
const productIndexName = "ProductIndex"

func catalogKey(ctx context.Context, catalogID string) *datastore.Key {
	return datastore.NewKey(ctx, catalogTypeName, catalogID, 0, nil)
}

func productKey(ctx context.Context, productID string) *datastore.Key {
	return datastore.NewKey(ctx, productTypeName, productID, 0, nil)
}
