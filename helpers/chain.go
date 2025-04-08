package helpers

import (
	"net/http"
)

func ChainMiddlewareHandlers(handlers ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for _, handler := range handlers {
			next = handler(next)
		}
		return next
	}
}

func ChainFuncs(handlers ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler(w, r)
		}
	}
}
