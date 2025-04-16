package handlers

import (
	"abysscore/common/errors/base"
	"abysscore/internal/adapters/controller/http/dto/request"
	"abysscore/internal/adapters/controller/http/dto/response"
	domainservice "abysscore/internal/domain/service"
	"abysscore/internal/infrastructure/metrics/tracer"
	"abysscore/pkg/validator"
	"context"
	"github.com/gofiber/fiber/v2"
)

type Authentication struct {
	authenticationService domainservice.AuthenticationService
	v                     *validator.Validator
	handlerName           string
}

func NewAuthenticationHandler(
	authenticationService domainservice.AuthenticationService,
	v *validator.Validator,
) *Authentication {
	return &Authentication{
		authenticationService: authenticationService,
		v:                     v,
		handlerName:           "authentication",
	}
}

func (a *Authentication) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var r request.AuthenticationRequest

	err := tracer.TraceFn(ctx, "validator.ValidateJSON", func(ctx context.Context) error {
		return validator.ValidateJSON(c, &r)
	})

	if err != nil {
		return err
	}

	var result *domainservice.AuthenticationResult

	result, err = tracer.TraceFnWithResult(ctx, "authenticationService.Register", func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
		return a.authenticationService.Register(ctx, r.ToCredentialsDTO())
	})

	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	return response.Success(c, result)
}

func (a *Authentication) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var r request.AuthenticationRequest

	err := tracer.TraceFn(ctx, "validator.ValidateJSON", func(ctx context.Context) error {
		return validator.ValidateJSON(c, &r)
	})

	var result *domainservice.AuthenticationResult

	result, err = tracer.TraceFnWithResult(ctx, "authenticationService.Authenticate", func(ctx context.Context) (*domainservice.AuthenticationResult, error) {
		return a.authenticationService.Authenticate(c.UserContext(), r.ToCredentialsDTO())
	})

	if err != nil {
		return base.ParseErrorOrInternalResponse(err, c)
	}

	return response.Success(c, result)
}
