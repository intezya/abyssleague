package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/dto/request"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
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

// @Router /api/auth/register [post].
func (h *AuthenticationHandler) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()

	req, err := getRequest[request.AuthenticationRequest](c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"authService.Register",
		func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
			return h.authenticationService.Register(ctx, req.ToCredentialsDTO())
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// @Router /api/auth/login [post].
func (h *AuthenticationHandler) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()

	req, err := getRequest[request.AuthenticationRequest](c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"authService.Authenticate",
		func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
			return h.authenticationService.Authenticate(ctx, req.ToCredentialsDTO())
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// @Router /api/auth/change_password [post].
func (h *AuthenticationHandler) ChangePassword(c *fiber.Ctx) error {
	ctx := c.UserContext()

	req, err := getRequest[request.PasswordChangeRequest](c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"authService.ChangePassword",
		func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
			return h.authenticationService.ChangePassword(ctx, req.ToDTO())
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}
