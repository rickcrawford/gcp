package handlers

import (
	"net/http"
)

const script = `
(function(){

	console.log('loaded');

})()
`

func scriptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(script))
	return
}
