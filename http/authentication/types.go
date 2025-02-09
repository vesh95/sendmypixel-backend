package authentication

import "errors"

var (
	// TypeNotAvailableError возникает при использования недопустимого типа аутентификации в заголовках
	TypeNotAvailableError = errors.New("invalid authentication type")

	// UserNotPresenceInStorageError возникает когда пользователь не найден в базе, но данные для аутентификации валидны
	UserNotPresenceInStorageError = errors.New("user is missing from the repository")

	// InvalidSignatureError возникает, когда переданные данные аутентификации валидны, но подпись не соответствует
	InvalidSignatureError = errors.New("invalid signature")

	// DataInvalidError возникает, когда переданные данные аутентификации не валидны
	DataInvalidError = errors.New("authentication data invalid")
)

type User struct {
	Id       int64
	Username string
}

type Authenticator interface {
	Authenticate(token string) (User, error)
}
