package applicationservice

import (
	"context"
	adaptererror "github.com/intezya/abyssleague/services/abysscore/internal/common/errors/adapter"
	applicationerror "github.com/intezya/abyssleague/services/abysscore/internal/common/errors/application"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/mailmessage"
	drivenports "github.com/intezya/abyssleague/services/abysscore/internal/domain/ports/driven"
	repositoryports "github.com/intezya/abyssleague/services/abysscore/internal/domain/repository"
)

type UserService struct {
	userRepository        repositoryports.UserRepository
	mailSender            drivenports.MailSender
	mailMessageRepository repositoryports.MailMessageRepository
}

func (s *UserService) SendCodeForEmailLink(
	ctx context.Context,
	user *dto.UserDTO,
	email string,
) error {
	if user.Email != nil {
		return applicationerror.ErrAccountAlreadyHasEmail
	}

	typedEmail, err := drivenports.NewEmail(email)

	if err != nil {
		return adaptererror.BadRequestFunc(err)
	}

	mailMessage := mailmessage.NewLinkEmailCodeMail(user.ID, email)

	s.mailMessageRepository.SaveLinkMailCodeData(ctx, mailMessage)
	err = s.mailSender.Send(ctx, typedEmail, mailMessage)

	if err != nil {
		return applicationerror.WrapServiceUnavailable(err)
	}

	return nil
}

func (s *UserService) EnterCodeForEmailLink(
	ctx context.Context,
	user *dto.UserDTO,
	verificationCode string,
) (
	*dto.UserDTO,
	error,
) {
	mailMessageData, err := s.mailMessageRepository.GetLinkMailCodeData(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	if mailMessageData.VerificationCode != verificationCode {
		return nil, applicationerror.ErrWrongVerificationCodeForEmailLink
	}

	result, err := s.userRepository.SetEmailIfNil(ctx, user.ID, mailMessageData.EmailForLink)

	if err != nil {
		return nil, err
	}

	return result, nil
}
