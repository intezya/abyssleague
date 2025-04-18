package server

import (
	"github.com/gofiber/fiber/v2"
)

// RouteGroup represents a group of routes with a common path prefix
type RouteGroup struct {
	prefix     string
	routes     []*RouteEntry
	middleware []fiber.Handler
}

// RouteEntry represents a single route within a group
type RouteEntry struct {
	path  string
	route *Route
}

// NewRouteGroup creates a new route group with the given prefix
func NewRouteGroup(prefix string) *RouteGroup {
	return &RouteGroup{
		prefix:     prefix,
		routes:     make([]*RouteEntry, 0),
		middleware: make([]fiber.Handler, 0),
	}
}

// Add adds a new route to the group
func (rg *RouteGroup) Add(relativePath string, route *Route) *RouteGroup {
	rg.routes = append(
		rg.routes, &RouteEntry{
			path:  relativePath,
			route: route,
		},
	)
	return rg
}

// Use adds middleware to the group
func (rg *RouteGroup) Use(middleware ...fiber.Handler) *RouteGroup {
	rg.middleware = append(rg.middleware, middleware...)
	return rg
}

// Register registers all routes in the group with the fiber app
func (rg *RouteGroup) Register(app *fiber.App, middlewareLinker *MiddlewareLinker) {
	group := app.Group(rg.prefix)

	// Apply group middleware
	for _, middleware := range rg.middleware {
		group.Use(middleware)
	}

	// Register all routes in the group
	for _, entry := range rg.routes {
		switch entry.route.Method {
		case MethodGet:
			group.Get(entry.path, middlewareLinker.buildMiddleware(entry.route))
		case MethodPost:
			group.Post(entry.path, middlewareLinker.buildMiddleware(entry.route))
		case MethodPut:
			group.Put(entry.path, middlewareLinker.buildMiddleware(entry.route))
		case MethodPatch:
			group.Patch(entry.path, middlewareLinker.buildMiddleware(entry.route))
		case MethodDelete:
			group.Delete(entry.path, middlewareLinker.buildMiddleware(entry.route))
		}
	}
}
