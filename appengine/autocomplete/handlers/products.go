package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"github.com/rickcrawford/gcp/appengine/autocomplete/common"
	"github.com/rickcrawford/gcp/appengine/autocomplete/managers"
	"github.com/rickcrawford/gcp/appengine/autocomplete/models"
)

type productCtxID int

const productCtxKey productCtxID = 0

type productHandler struct {
	manager managers.ProductManager
}

func (productHandler) get(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(common.Response{
		Data: r.Context().Value(productCtxKey),
		Metadata: map[string]interface{}{
			"status": http.StatusOK,
		},
	})
}

func (h productHandler) saveBatch(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		common.WriteBadRequest(w)
		return
	}

	ctx := r.Context()
	catalog := ctx.Value(catalogCtxKey).(*models.Catalog)

	var products []models.Product
	err := json.NewDecoder(r.Body).Decode(&products)
	r.Body.Close()
	if err != nil {
		common.WriteError(w, err)
		return
	}

	err = h.manager.SaveAll(ctx, catalog.CatalogID, products)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	json.NewEncoder(w).Encode(common.Response{
		Metadata: map[string]interface{}{
			"status": http.StatusOK,
			"count":  len(products),
		},
	})
}

func (h productHandler) save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	product := &models.Product{}
	h.update(w, r.WithContext(context.WithValue(ctx, productCtxKey, product)))
}

func (h productHandler) update(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		common.WriteBadRequest(w)
		return
	}

	ctx := r.Context()
	product := ctx.Value(productCtxKey).(*models.Product)
	catalog := ctx.Value(catalogCtxKey).(*models.Catalog)

	err := json.NewDecoder(r.Body).Decode(product)
	r.Body.Close()
	if err != nil {
		common.WriteError(w, err)
		return
	}

	err = h.manager.Save(ctx, catalog.CatalogID, product)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	json.NewEncoder(w).Encode(common.Response{
		Data: product,
		Metadata: map[string]interface{}{
			"status": http.StatusOK,
		},
	})
}

func (h productHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	product := ctx.Value(productCtxKey).(*models.Product)
	catalog := ctx.Value(catalogCtxKey).(*models.Catalog)

	err := h.manager.Delete(ctx, catalog.CatalogID, product.ProductID)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	json.NewEncoder(w).Encode(common.Response{
		Metadata: map[string]interface{}{
			"status": http.StatusOK,
		},
	})

}

func (h productHandler) context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		catalog := ctx.Value(catalogCtxKey).(*models.Catalog)
		productID := chi.URLParam(r, "productID")

		product, _ := h.manager.Get(ctx, catalog.CatalogID, productID)
		if product == nil {
			common.WriteNotFound(w)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, productCtxKey, product)))
	})
}
