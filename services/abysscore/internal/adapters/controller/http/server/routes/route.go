package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/schema/access_level"
)

type (
	MatchRequirement int
	RateLimit        int
	Method           string
)

const (
	MatchIrrelevant MatchRequirement = iota
	MustBeInMatch
	MustNotBeInMatch
)

const (
	DisableRateLimit RateLimit = iota
	AuthRateLimit
	DefaultRateLimit
)

const (
	MethodGet    = fiber.MethodGet
	MethodPost   = fiber.MethodPost
	MethodPut    = fiber.MethodPut
	MethodPatch  = fiber.MethodPatch
	MethodDelete = fiber.MethodDelete
)

type Route struct {
	Handler fiber.Handler
	Method
	RequireAuthentication bool
	AccessLevel           *access_level.AccessLevel
	MatchRequirement      MatchRequirement
	RateLimit             RateLimit
}

type RouteOption func(*Route)

func NewRoute(handler fiber.Handler, method Method, opts ...RouteOption) *Route {
	route := &Route{
		Handler:               handler,
		Method:                method,
		RequireAuthentication: true,
		AccessLevel:           nil,
		MatchRequirement:      MatchIrrelevant,
		RateLimit:             DefaultRateLimit,
	}

	for _, opt := range opts {
		opt(route)
	}

	return route
}

func WithAccessLevel(level access_level.AccessLevel) RouteOption {
	return func(r *Route) {
		r.AccessLevel = &level
	}
}

func WithMatchRequirement(req MatchRequirement) RouteOption {
	return func(r *Route) {
		r.MatchRequirement = req
	}
}

func WithRateLimit(rate RateLimit) RouteOption {
	return func(r *Route) {
		r.RateLimit = rate
	}
}
