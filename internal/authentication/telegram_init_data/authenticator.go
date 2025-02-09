package telegram_init_data

import (
	"backend/http/authentication"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

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

func (t *TelegramMiniApp) Authenticate(token string) (authentication.User, error) {
	parsedInitData, err := initdata.Parse(token)
	if err != nil {
		return authentication.User{}, authentication.DataInvalidError
	}

	err = initdata.Validate(token, t.secret, 0)
	if err != nil {
		return authentication.User{}, authentication.InvalidSignatureError
	}

	user, err := t.authStorage.Upsert(parsedInitData.User)
	if err != nil {
		return authentication.User{}, err
	}

	return user, nil
}
