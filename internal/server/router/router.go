// Package router provides functionality to set up and run the HTTP server,
// manage routes, and apply middleware based on access levels.
//
// It configures and initializes routes with varying access requirements:
// - Public routes: Accessible without authentication.
// - Protected routes: Require authentication.
package router

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	api "github.com/mohit/finance-go/internal/server/rest"
)

// Router manages the HTTP server and its dependencies,
// including controllers and JWT authentication.
type Router struct {
	addr                     string
	baseURL                  string
	restfullControllers      []api.IController
	graphQlController        http.Handler
	authorizationMiddleware  func(http.Handler) http.Handler
	populateClaimsMiddleware func(http.Handler) http.Handler
	rateLimitMiddleware      func(http.Handler) http.Handler
}

// Config holds configuration settings for creating a new Router instance.
type Config struct {
	Addr                     string            // Address to listen on
	BaseURL                  string            // Base URL for API routes
	RestfullControllers      []api.IController // List of controllers
	GraphQlController        http.Handler
	AuthorizationMiddleware  func(http.Handler) http.Handler
	PopulateClaimsMiddleware func(http.Handler) http.Handler
	RateLimitMiddleware      func(http.Handler) http.Handler
}

// NewRouter creates a new Router instance with the given configuration.
// It initializes the router with address, base URL, controllers, and JWT service.
func NewRouter(config Config) *Router {
	return &Router{
		addr:                     config.Addr,
		baseURL:                  config.BaseURL,
		restfullControllers:      config.RestfullControllers,
		graphQlController:        config.GraphQlController,
		authorizationMiddleware:  config.AuthorizationMiddleware,
		populateClaimsMiddleware: config.PopulateClaimsMiddleware,
		rateLimitMiddleware:      config.RateLimitMiddleware,
	}
}

// Run starts the HTTP server and sets up routes with different access levels.
//
// Routes are grouped and managed under the base URL, with the following access levels:
// - Public routes: No authentication required.
// - Protected routes: Authentication required.
func (r *Router) Run() error {
	mux := http.NewServeMux()

	for _, c := range r.restfullControllers {
		c.RegisterPublic(mux)
		c.RegisterProtected(mux, r.authorizationMiddleware)
	}

	mux.Handle("/api/graph/", r.populateClaimsMiddleware(playground.Handler("GraphQL playground", "/api/graph/query")))
	mux.Handle("/api/graph/query", r.populateClaimsMiddleware(r.graphQlController))

	// Serve the frontend static files
	mux.Handle("/", http.FileServer(http.Dir("./frontend")))

	// Wrap the entire mux with the rate limit middleware
	var handler http.Handler = mux
	if r.rateLimitMiddleware != nil {
		handler = r.rateLimitMiddleware(mux)
	}

	log.Println("Listening on", r.addr)
	return http.ListenAndServe(r.addr, handler)
}
