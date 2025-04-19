package routes

import (
	"github.com/gofiber/fiber/v2"
)

// RouteGroup represents a group of routes with a common path prefix
type RouteGroup struct {
	Prefix     string
	Routes     []*RouteEntry
	middleware []fiber.Handler
}

// RouteEntry represents a single route within a group
type RouteEntry struct {
	Path  string
	Route *Route
}

// NewRouteGroup creates a new route group with the given prefix
func NewRouteGroup(prefix string) *RouteGroup {
	return &RouteGroup{
		Prefix:     prefix,
		Routes:     make([]*RouteEntry, 0),
		middleware: make([]fiber.Handler, 0),
	}
}

// Add adds a new route to the group
func (rg *RouteGroup) Add(relativePath string, route *Route) *RouteGroup {
	rg.Routes = append(
		rg.Routes, &RouteEntry{
			Path:  relativePath,
			Route: route,
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
	group := app.Group(rg.Prefix)

	// Apply group middleware
	for _, middleware := range rg.middleware {
		group.Use(middleware)
	}

	// Register all routes in the group
	for _, entry := range rg.Routes {
		handlers := []fiber.Handler{
			middlewareLinker.loggingMiddleware.Handle(),
			middlewareLinker.recoverMiddleware.Handle(),
		}

		// Rate limiting middleware
		switch entry.Route.RateLimit {
		case DefaultRateLimit:
			handlers = append(handlers, middlewareLinker.rateLimitMiddleware.HandleDefault())
		case AuthRateLimit:
			handlers = append(handlers, middlewareLinker.rateLimitMiddleware.HandleForAuth())
		default:
		}

		// Authentication middleware
		if entry.Route.RequireAuthentication {
			handlers = append(handlers, middlewareLinker.authenticationMiddleware.Handle())

			if entry.Route.AccessLevel != nil {
				handlers = append(handlers, createAccessLevelChecker(entry.Route.AccessLevel))
			}

			if entry.Route.MatchRequirement != MatchIrrelevant {
				handlers = append(handlers, createMatchRequirementChecker(entry.Route.MatchRequirement))
			}
		}

		handlers = append(handlers, entry.Route.Handler)

		switch entry.Route.Method {
		case MethodGet:
			group.Get(entry.Path, handlers...)
		case MethodPost:
			group.Post(entry.Path, handlers...)
		case MethodPut:
			group.Put(entry.Path, handlers...)
		case MethodPatch:
			group.Patch(entry.Path, handlers...)
		case MethodDelete:
			group.Delete(entry.Path, handlers...)
		}
	}
}
