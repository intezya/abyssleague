package applicationservice

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/mailmessage"
	drivenports "github.com/intezya/abyssleague/services/abysscore/internal/domain/ports/driven"
	repositoryports "github.com/intezya/abyssleague/services/abysscore/internal/domain/repository"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
)

type AccountService struct {
	userRepository        repositoryports.UserRepository
	mailSender            drivenports.MailSender
	mailMessageRepository repositoryports.MailMessageRepository
}

func NewAccountService(
	userRepository repositoryports.UserRepository,
	mailSender drivenports.MailSender,
	mailMessageRepository repositoryports.MailMessageRepository,
) *AccountService {
	return &AccountService{
		userRepository:        userRepository,
		mailSender:            mailSender,
		mailMessageRepository: mailMessageRepository,
	}
}

func (s *AccountService) SendCodeForEmailLink(
	ctx context.Context,
	user *dto.UserDTO,
	email string,
) error {
	ctx, span := tracer.StartSpan(ctx, "AccountService.SendCodeForEmailLink")
	defer span.End()

	const newLinkEmailCodeExpireMinutes = 5

	if user.Email != nil {
		return apperrors.ErrAccountAlreadyHasEmail
	}

	typedEmail, err := drivenports.NewEmail(email)
	if err != nil {
		return apperrors.WrapBadRequest(err)
	}

	exists := s.userRepository.ExistsByEmail(ctx, email)

	if exists {
		return apperrors.ErrEmailConflict
	}

	sentMailMessage, err := s.mailMessageRepository.GetLinkMailCodeData(ctx, user.ID)

	if err == nil && sentMailMessage != nil &&
		sentMailMessage.CreatedAt.After(
			time.Now().Add(-time.Minute*1),
		) { // TODO: move to "if expired" func (move timeout to const)
		return apperrors.TooManyEmailLinkRequests // if already sent
	}

	mailMessage := mailmessage.NewLinkEmailCodeMail(user.ID, email, newLinkEmailCodeExpireMinutes)

	err = s.mailSender.Send(ctx, mailMessage.Message, typedEmail.String())
	if err != nil {
		return apperrors.WrapServiceUnavailable(err)
	}

	go s.mailMessageRepository.SaveLinkMailCodeData(ctx, mailMessage, newLinkEmailCodeExpireMinutes)

	return nil
}

func (s *AccountService) EnterCodeForEmailLink(
	ctx context.Context,
	user *dto.UserDTO,
	verificationCode string,
) (
	*dto.UserDTO,
	error,
) {
	ctx, span := tracer.StartSpan(ctx, "AccountService.EnterCodeForEmailLink")
	defer span.End()

	mailMessageData, err := s.mailMessageRepository.GetLinkMailCodeData(ctx, user.ID)
	if err != nil {
		return nil, apperrors.WrapWrongVerificationCodeForEmailLink(err)
	}

	if mailMessageData.VerificationCode != verificationCode {
		return nil, apperrors.ErrWrongVerificationCodeForEmailLink
	}

	result, err := s.userRepository.SetEmailIfNil(ctx, user.ID, mailMessageData.EmailForLink)
	if err != nil {
		return nil, err
	}

	return result, nil
}
