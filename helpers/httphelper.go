package helpers

import "net/http"

func SetHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, public")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("ETag", "456")
}
