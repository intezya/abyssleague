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
// @Summary Register a new user
// @Description Creates a new user account with the provided credentials
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.AuthenticationRequest true "User registration details"
// @Success 200 {object} domainservice.AuthenticationResult "User successfully registered"
// @Failure 409 {object} examples.UsernameConflictResponse "Conflict - user with this username already exists"
// @Failure 409 {object} examples.HardwareIDConflictResponse "Conflict - only one account per device allowed"
// @Failure 422 {object} examples.UnprocessableErrorResponse "Unprocessable entity - validation error"
// @Failure 429 {object} examples.TooManyRequestsResponse "Too many requests - received too many auth requests"
// @Failure 500 {object} base.ErrorResponse "Internal server error"
// @Router /api/auth/register [post]
func (a *AuthenticationHandler) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()
	r := &request.AuthenticationRequest{}

	if err := a.validateRequest(r, c); err != nil {
		return a.handleError(err, c)
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
// @Summary Authenticate user
// @Description Authenticates a user with username and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.AuthenticationRequest true "User login credentials"
// @Success 200 {object} domainservice.AuthenticationResult "Successfully authenticated"
// @Failure 401 {object} examples.UserWrongPasswordResponse "Unauthorized - wrong password"
// @Failure 401 {object} examples.UserWrongHardwareIDResponse "Unauthorized - wrong hardware id"
// @Failure 404 {object} examples.UserNotFoundResponse "Not found - user with this username not found"
// @Failure 409 {object} examples.UsernameConflictResponse "Conflict - user with this username already exists"
// @Failure 409 {object} examples.HardwareIDConflictResponse "Conflict - only one account per device allowed"
// @Failure 422 {object} examples.UnprocessableErrorResponse "Unprocessable entity - validation error"
// @Failure 429 {object} examples.TooManyRequestsResponse "Too many requests - received too many auth requests"
// @Router /api/auth/login [post]
func (a *AuthenticationHandler) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()

	r := &request.AuthenticationRequest{}
	if err := a.validateRequest(r, c); err != nil {
		return a.handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(ctx, "authService.Authenticate", func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
		return a.authenticationService.Authenticate(ctx, r.ToCredentialsDTO())
	})

	if err != nil {
		return a.handleError(err, c)
	}

	return a.sendSuccess(result, c)
}

// ChangePassword handles changing user password
// @Summary Change user password
// @Description Changes the password for an existing user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.PasswordChangeRequest true "Password change details"
// @Success 200 {object} domainservice.AuthenticationResult "Password successfully changed"
// @Failure 401 {object} examples.UserWrongPasswordResponse "Unauthorized - wrong password"
// @Failure 404 {object} examples.UserNotFoundResponse "Not found - user with this username not found"
// @Failure 422 {object} examples.UnprocessableErrorResponse "Unprocessable entity - validation error"
// @Failure 429 {object} examples.TooManyRequestsResponse "Too many requests - received too many auth requests"
// @Router /api/auth/change_password [post]
func (a *AuthenticationHandler) ChangePassword(c *fiber.Ctx) error {
	ctx := c.UserContext()

	r := &request.PasswordChangeRequest{}
	if err := a.validateRequest(r, c); err != nil {
		return a.handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(ctx, "authService.ChangePassword", func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
		return a.authenticationService.ChangePassword(ctx, r.ToDTO())
	})

	if err != nil {
		return a.handleError(err, c)
	}

	return a.sendSuccess(result, c)
}
