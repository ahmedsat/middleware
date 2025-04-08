package handlers

import (
	"net/http"
)

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.webp")
}
