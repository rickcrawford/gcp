package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rickcrawford/gcp/appengine/autocomplete/common"
	"github.com/rickcrawford/gcp/appengine/autocomplete/managers"
	"github.com/rickcrawford/gcp/appengine/autocomplete/models"
)

type searchHandler struct {
	manager managers.ProductManager
}

func (h searchHandler) search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if strings.EqualFold(r.Method, "HEAD") {
		return
	}

	ctx := r.Context()
	catalog := ctx.Value(catalogCtxKey).(*models.Catalog)

	var etag string
	if catalog.Batch != "" {
		etag = fmt.Sprintf(`W/"%s"`, catalog.Batch)
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	limit := 10
	if r.FormValue("limit") != "" {
		limitVal, _ := strconv.Atoi(r.FormValue("limit"))
		if limitVal > 0 {
			limit = limitVal
		}
	}

	cache := 0
	if r.FormValue("cache") != "" {
		cacheVal, _ := strconv.Atoi(r.FormValue("cache"))
		if cacheVal > 0 {
			cache = cacheVal
		}
	}

	var query string
	if r.FormValue("prefix") != "" {
		prefix := common.FormatPrefix(r.FormValue("prefix"), "_")
		query = fmt.Sprintf("prefix:%s", prefix)
	} else if r.FormValue("query") != "" {
		query = r.FormValue("query")
	} else {
		common.WriteBadRequest(w)
		return
	}

	products, err := h.manager.Search(ctx, catalog.CatalogID, models.SearchQuery{
		Query: query,
		Limit: limit,
	})
	if err != nil {
		common.WriteError(w, err)
		return
	}

	if cache > 0 {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d", cache))
	}

	if etag != "" {
		w.Header().Add("ETag", etag)
	}

	json.NewEncoder(w).Encode(common.Response{
		Data: products,
		Metadata: map[string]interface{}{
			"status": http.StatusOK,
			"count":  len(products),
			"query":  query,
			"cache":  cache,
		},
	})
}
