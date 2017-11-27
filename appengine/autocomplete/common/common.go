package common

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Response is a response struct for results
type Response struct {
	Data     interface{}            `json:"data,omitempty"`
	Errors   []string               `json:"errors,omitempty"`
	Metadata map[string]interface{} `json:"metadata"`
}

// MaxPrefixLength is the maximum prefix length
const MaxPrefixLength = 10

var NullTime = time.Time{}

var pattern = regexp.MustCompile(`[^\p{L}\p{N}]+`)

func WriteNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Metadata: map[string]interface{}{
			"status": http.StatusNotFound,
		},
	})
}

func WriteBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{
		Metadata: map[string]interface{}{
			"status": http.StatusBadRequest,
		},
	})
}

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{
		Errors: []string{err.Error()},
		Metadata: map[string]interface{}{
			"status": http.StatusBadRequest,
		},
	})
}

func FormatPrefix(prefix, replacement string) string {
	return pattern.ReplaceAllString(strings.ToLower(prefix), replacement)
}

func SliceContains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
