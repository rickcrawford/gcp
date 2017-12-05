package models

// Response is a response struct for results
type Response struct {
	Data     interface{}            `json:"data,omitempty"`
	Errors   []string               `json:"errors,omitempty"`
	Metadata map[string]interface{} `json:"metadata"`
}
