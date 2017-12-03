package managers

import (
	"context"
	"strconv"

	"google.golang.org/appengine/datastore"
)

const productTypeName = "Product"
const productIndexName = "ProductIndex"

func productKey(ctx context.Context, productID int) *datastore.Key {
	return datastore.NewKey(ctx, productTypeName, strconv.Itoa(productID), 0, nil)
}
