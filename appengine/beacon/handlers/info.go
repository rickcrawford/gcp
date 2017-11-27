package handlers

import (
	"bytes"
	"net/http"
	"os"
	"time"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Beacon-Debug") != "true" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", time.Time{}.Format(time.RFC822))
	w.WriteHeader(http.StatusOK)

	var buf bytes.Buffer
	addValue(&buf, "RemoteAddress", r.RemoteAddr)
	addValue(&buf, "Method", r.Method)
	addValue(&buf, "Host", r.Host)
	addValue(&buf, "URL", r.URL.String())
	buf.WriteString("-headers--------------------------------\n")
	for key := range r.Header {
		addValue(&buf, key, r.Header.Get(key))
	}

	buf.WriteString("-env------------------------------------\n")
	for _, key := range os.Environ() {
		addValue(&buf, key, os.Getenv(key))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(buf.Bytes())
}

func addValue(buf *bytes.Buffer, key, value string) {
	buf.WriteString(key)
	buf.WriteByte(':')
	buf.WriteString(value)
	buf.WriteByte('\n')
}
