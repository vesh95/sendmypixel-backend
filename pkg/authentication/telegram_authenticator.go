package authentication

import (
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

const UserContextKey = "user"

type TelegramMiniApp struct {
	secret      string
	authStorage TelegramAuthenticationStorage
}

func NewTelegramMiniAppAuthenticator(secret string, authStorage TelegramAuthenticationStorage) *TelegramMiniApp {
	return &TelegramMiniApp{
		secret:      secret,
		authStorage: authStorage,
	}
}

func (t *TelegramMiniApp) Authenticate(token string) (User, error) {
	parsedInitData, err := initdata.Parse(token)
	if err != nil {
		return User{}, DataInvalidError
	}

	err = initdata.Validate(token, t.secret, 0)
	if err != nil {
		return User{}, InvalidSignatureError
	}

	user, err := t.authStorage.Upsert(parsedInitData.User)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

type TelegramAuthenticationStorage interface {
	Upsert(user initdata.User) (User, error)
}
