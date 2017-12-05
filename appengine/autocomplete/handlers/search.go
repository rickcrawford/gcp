package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rickcrawford/gcp/appengine/autocomplete/common"
	"github.com/rickcrawford/gcp/appengine/autocomplete/managers"
	"github.com/rickcrawford/gcp/common/models"
)

type searchHandler struct {
	manager managers.ProductManager
}

func (h searchHandler) search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if strings.EqualFold(r.Method, "HEAD") {
		return
	}

	ctx := r.Context()

	limit := 10
	if r.FormValue("limit") != "" {
		limitVal, _ := strconv.Atoi(r.FormValue("limit"))
		if limitVal > 0 {
			limit = limitVal
		}
	}

	var query string
	if r.FormValue("q") != "" {
		prefix := common.FormatPrefix(r.FormValue("q"), "_")
		if len(prefix) > managers.MaxPrefixLength {
			prefix = prefix[:managers.MaxPrefixLength]
		}
		query = fmt.Sprintf("prefix:%s", prefix)
	} else if r.FormValue("q") != "" {
		query = r.FormValue("q")
	} else {
		common.WriteBadRequest(w)
		return
	}

	products, err := h.manager.Search(ctx, managers.SearchQuery{
		Query: query,
		Limit: limit,
	})
	if err != nil {
		common.WriteError(w, err)
		return
	}
	w.Header().Add("Cache-Control", "max-age=86400")

	json.NewEncoder(w).Encode(models.Response{
		Data: products,
		Metadata: map[string]interface{}{
			"status": http.StatusOK,
			"count":  len(products),
			"query":  query,
		},
	})
}
