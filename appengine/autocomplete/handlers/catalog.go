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

type catalogCtxID int

const catalogCtxKey catalogCtxID = 0

type catalogHandler struct {
	manager managers.CatalogManager
}

func (h catalogHandler) get(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(common.Response{
		Data: r.Context().Value(catalogCtxKey),
		Metadata: map[string]interface{}{
			"status": http.StatusOK,
		},
	})
}

func (h catalogHandler) getBatch(w http.ResponseWriter, r *http.Request) {
	catalog := r.Context().Value(catalogCtxKey).(*models.Catalog)
	json.NewEncoder(w).Encode(common.Response{
		Data: catalog.Batch,
		Metadata: map[string]interface{}{
			"status": http.StatusOK,
		},
	})
}

func (h catalogHandler) save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.update(w, r.WithContext(context.WithValue(ctx, catalogCtxKey, new(models.Catalog))))
}

func (h catalogHandler) update(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		common.WriteBadRequest(w)
		return
	}

	ctx := r.Context()
	catalog := ctx.Value(catalogCtxKey).(*models.Catalog)

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(catalog)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	err = h.manager.Save(ctx, catalog)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	json.NewEncoder(w).Encode(common.Response{
		Data: catalog,
		Metadata: map[string]interface{}{
			"status": http.StatusOK,
		},
	})
}

func (h catalogHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	catalog := ctx.Value(catalogCtxKey).(*models.Catalog)

	err := h.manager.Delete(ctx, catalog.CatalogID)
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

func (h catalogHandler) context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		catalogID := chi.URLParam(r, "catalogID")
		ctx := r.Context()

		catalog, _ := h.manager.Get(ctx, catalogID)
		if catalog == nil {
			common.WriteNotFound(w)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, catalogCtxKey, catalog)))
	})
}
