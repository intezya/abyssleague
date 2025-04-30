package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/dto/request"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
)

type AccountHandler struct {
	accountService domainservice.AccountService
}

func NewAccountHandler(accountService domainservice.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// SendCodeForEmailLink handles sending verification code for email linking
//
//	@Summary		Send email verification code
//	@Description	Sends verification code for email linking
//	@Tags			Account
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	request.LinkEmailRequest	true	"Email for linking"
//	@Success		204		"Code successfully sent"
//	@Failure		400		{object}	examples.BadRequestResponse				"Bad request - missed request fields"
//	@Failure		409		{object}	examples.AccountAlreadyHasLinkedEmail	"Conflict - user already has linked email"
//	@Failure		409		{object}	examples.EmailConflict					"Conflict - someone already has this email as linked"
//	@Failure		422		{object}	examples.UnprocessableEntityResponse	"Unprocessable entity - invalid request types"
//	@Failure		429		{object}	examples.TooManyRequestsResponse		"Too many requests - received too many requests"
//	@Router			/api/account/email/get_code [post].
func (h *AccountHandler) SendCodeForEmailLink(c *fiber.Ctx) error {
	ctx := c.UserContext()

	user, err := extractUser(ctx)
	if err != nil {
		return handleError(err, c)
	}

	req, err := getAndValidateRequest[request.LinkEmailRequest](c)
	if err != nil {
		return handleError(err, c)
	}

	err = h.accountService.SendCodeForEmailLink(ctx, user, req.Email)

	if err != nil {
		return handleError(err, c)
	}

	return sendNoContent(c)
}

// EnterCodeForEmailLink verifies email verification code and links email to account
//
//	@Summary		Verify sent code
//	@Description	Verifies sent code and links email to account
//	@Tags			Account
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.EnterCodeForEmailLinkRequest	true	"Verification code from email"
//	@Success		200		{object}	dto.UserDTO								"Email successfully linked"
//	@Failure		400		{object}	examples.BadRequestResponse				"Bad request - missed request fields"
//	@Failure		400		{object}	examples.BadRequestResponse				"Bad request - wrong verification code"
//	@Failure		409		{object}	examples.AccountAlreadyHasLinkedEmail	"Conflict - user already has linked email"
//	@Failure		422		{object}	examples.UnprocessableEntityResponse	"Unprocessable entity - invalid request types"
//	@Failure		429		{object}	examples.TooManyRequestsResponse		"Too many requests - received too many requests"
//	@Router			/api/account/email/enter_code [post].
func (h *AccountHandler) EnterCodeForEmailLink(c *fiber.Ctx) error {
	ctx := c.UserContext()

	user, err := extractUser(ctx)
	if err != nil {
		return handleError(err, c)
	}

	req, err := getAndValidateRequest[request.EnterCodeForEmailLinkRequest](c)
	if err != nil {
		return handleError(err, c)
	}

	result, err := h.accountService.EnterCodeForEmailLink(ctx, user, req.VerificationCode)

	if err != nil {
		return handleError(err, c)
	}

	return sendSuccess(result, c)
}
