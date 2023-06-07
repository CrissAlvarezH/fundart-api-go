package ports

type MessageProvider string

const (
	MessageProviderSMS   MessageProvider = "sms"
	MessageProviderEmail MessageProvider = "email"
)

type VerificationCodeManager interface {
	Send(code string, target string, provider MessageProvider) error
}
