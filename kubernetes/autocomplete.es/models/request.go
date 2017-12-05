package models

import "time"

// RequestData used by cloud function for requests
type RequestData struct {
	RawHeaders  []string          `json:"rawHeaders"`
	Headers     map[string]string `json:"headers"`
	OriginalURL string            `json:"originalUrl"`
	URL         string            `json:"url"`
	Query       string            `json:"query"`
	Path        string            `json:"query"`
	Params      map[string]string `json:"params"`
	Method      string            `json:"method"`
	HTTPVersion string            `json:"httpVersion"`
	Hostname    string            `json:"hostname"`
	Secure      bool              `json:"secure"`
	Protocol    string            `json:"protocol"`
	Cookies     map[string]string `json:"cookies"`
	IP          string            `json:"ip"`
	IPs         []string          `json:"ips"`
	Data        string            `json:"data"`
	Time        struct {
		Request time.Time `json:"request"`
		Current time.Time `json:"current"`
		Skew    int       `json:"skew"`
	} `json:"time"`
	Token string `json:"token"`
}
