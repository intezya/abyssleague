package handlers

import (
	"abysscore/internal/adapters/controller/http/dto/request"
	"abysscore/internal/adapters/controller/http/dto/response"
	"abysscore/internal/common/errors/base"
	domainservice "abysscore/internal/domain/service"
	"abysscore/internal/infrastructure/metrics/tracer"
	"abysscore/internal/pkg/validator"
	"context"
	"github.com/gofiber/fiber/v2"
)

type AuthenticationHandler struct {
	authenticationService domainservice.AuthenticationService
}

func NewAuthenticationHandler(
	authenticationService domainservice.AuthenticationService,
) *AuthenticationHandler {
	return &AuthenticationHandler{
		authenticationService: authenticationService,
	}
}

// Register handles user registration
func (a *AuthenticationHandler) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var r = &request.AuthenticationRequest{}
	err := tracer.TraceFn(ctx, "validator.ValidateJSON", func(ctx context.Context) error {
		return validator.ValidateJSON(r, c)
	})
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	var result *domainservice.AuthenticationResult

	result, err = tracer.TraceFnWithResult(ctx, "authenticationService.Register", func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
		return a.authenticationService.Register(ctx, r.ToCredentialsDTO())
	})

	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	return response.Success(result, c)
}

// Login handles user authentication
func (a *AuthenticationHandler) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var r = &request.AuthenticationRequest{}

	err := tracer.TraceFn(ctx, "validator.ValidateJSON", func(ctx context.Context) error {
		return validator.ValidateJSON(r, c)
	})
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	var result *domainservice.AuthenticationResult

	result, err = tracer.TraceFnWithResult(ctx, "authenticationService.Authenticate", func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
		return a.authenticationService.Authenticate(ctx, r.ToCredentialsDTO())
	})
	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	return response.Success(result, c)
}
