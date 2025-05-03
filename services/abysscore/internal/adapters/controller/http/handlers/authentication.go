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

// Register handles user registration
//
//	@Summary		Register a new user
//	@Description	Creates a new user account with the provided credentials
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.AuthenticationRequest			true	"User registration details"
//	@Success		200		{object}	examples.AuthenticationSuccessResponse	"User successfully registered"
//	@Failure		400		{object}	examples.BadRequestResponse				"Bad request - missed request fields"
//	@Failure		409		{object}	examples.UsernameConflictResponse		"Conflict - user with this username already exists"
//	@Failure		409		{object}	examples.HardwareIDConflictResponse		"Conflict - only one account per device allowed"
//	@Failure		422		{object}	examples.UnprocessableEntityResponse	"Unprocessable entity - invalid request types"
//	@Failure		429		{object}	examples.TooManyRequestsResponse		"Too many requests - received too many auth requests"
//	@Router			/api/auth/register [post].
func (h *AuthenticationHandler) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()

	req, err := getAndValidateRequest[request.AuthenticationRequest](c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"authenticationService.Register",
		func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
			return h.authenticationService.Register(ctx, req.ToCredentialsDTO())
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// Login handles user authentication
//
//	@Summary		Authenticate user
//	@Description	Authenticates a user with username and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.AuthenticationRequest			true	"User login credentials"
//	@Success		200		{object}	examples.AuthenticationSuccessResponse	"Successfully authenticated"
//	@Failure		400		{object}	examples.BadRequestResponse				"Bad request - missed request fields"
//	@Failure		401		{object}	examples.UserWrongPasswordResponse		"Unauthorized - wrong password"
//	@Failure		401		{object}	examples.UserWrongHardwareIDResponse	"Unauthorized - wrong hardware id"
//	@Failure		404		{object}	examples.UserNotFoundResponse			"Not found - user with this username not found"
//	@Failure		409		{object}	examples.UsernameConflictResponse		"Conflict - user with this username already exists"
//	@Failure		409		{object}	examples.HardwareIDConflictResponse		"Conflict - only one account per device allowed"
//	@Failure		422		{object}	examples.UnprocessableEntityResponse	"Unprocessable entity - invalid request types"
//	@Failure		429		{object}	examples.TooManyRequestsResponse		"Too many requests - received too many auth requests"
//	@Router			/api/auth/login [post].
func (h *AuthenticationHandler) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()

	req, err := getAndValidateRequest[request.AuthenticationRequest](c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"authenticationService.Authenticate",
		func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
			return h.authenticationService.Authenticate(ctx, req.ToCredentialsDTO())
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}

// ChangePassword handles changing user password
//
//	@Summary		Change user password
//	@Description	Changes the password for an existing user
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.PasswordChangeRequest			true	"Password change details"
//	@Success		200		{object}	examples.AuthenticationSuccessResponse	"Password successfully changed"
//	@Failure		400		{object}	examples.BadRequestResponse				"Bad request - missed request fields"
//	@Failure		401		{object}	examples.UserWrongPasswordResponse		"Unauthorized - wrong password"
//	@Failure		404		{object}	examples.UserNotFoundResponse			"Not found - user with this username not found"
//	@Failure		422		{object}	examples.UnprocessableEntityResponse	"Unprocessable entity - invalid request types"
//	@Failure		429		{object}	examples.TooManyRequestsResponse		"Too many requests - received too many auth requests"
//	@Router			/api/auth/change_password [post].
func (h *AuthenticationHandler) ChangePassword(c *fiber.Ctx) error {
	ctx := c.UserContext()

	req, err := getAndValidateRequest[request.PasswordChangeRequest](c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := tracer.TraceFnWithResult(
		ctx,
		"authenticationService.ChangePassword",
		func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
			return h.authenticationService.ChangePassword(ctx, req.ToDTO())
		},
	)
	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}
