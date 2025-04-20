package handlers

import (
	"abysscore/internal/adapters/controller/http/dto/request"
	domainservice "abysscore/internal/domain/service"
	"abysscore/internal/infrastructure/metrics/tracer"
	"context"
	"github.com/gofiber/fiber/v2"
)

type AuthenticationHandler struct {
	BaseHandler

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
	r := &request.AuthenticationRequest{}

	if err := a.validateRequest(r, c); err != nil {
		return err
	}

	result, err := tracer.TraceFnWithResult(ctx, "authService.Register", func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
		return a.authenticationService.Register(ctx, r.ToCredentialsDTO())
	})

	if err != nil {
		return a.handleError(err, c)
	}

	return a.sendSuccess(result, c)
}

// Login handles user authentication
func (a *AuthenticationHandler) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()
	r := &request.AuthenticationRequest{}

	if err := a.validateRequest(r, c); err != nil {
		return err
	}

	result, err := tracer.TraceFnWithResult(ctx, "authService.Authenticate", func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
		return a.authenticationService.Authenticate(ctx, r.ToCredentialsDTO())
	})

	if err != nil {
		return a.handleError(err, c)
	}

	return a.sendSuccess(result, c)
}

func (a *AuthenticationHandler) ChangePassword(c *fiber.Ctx) error {
	ctx := c.UserContext()
	r := &request.PasswordChangeRequest{}

	if err := a.validateRequest(r, c); err != nil {
		return err
	}

	result, err := tracer.TraceFnWithResult(ctx, "authService.ChangePassword", func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
		return a.authenticationService.ChangePassword(ctx, r.ToDTO())
	})

	if err != nil {
		return a.handleError(err, c)
	}

	return a.sendSuccess(result, c)
}
