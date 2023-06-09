package ports

import "errors"

var (
	InvalidValidationCode = errors.New("invalid validation code")
)

type MessageProvider string

type VerificationCodeManager interface {
	SendEmailToVerifyAccount(code string, email string) error
	SendEmailToRecoverPassword(code string, email string) error
}
