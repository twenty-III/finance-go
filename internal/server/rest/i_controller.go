// Package api defines an interface for controllers to register routes
// with different access levels.
package api

import "net/http"

// IController outlines methods for route registration:
// - Public: No authentication required.
// - Protected: Requires authentication.
type IController interface {
	// RegisterPublic sets up public routes.
	RegisterPublic(mux *http.ServeMux)

	// RegisterProtected sets up routes that require authentication.
	RegisterProtected(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler)
}
