package telegram_init_data

import (
	"backend/http/authentication"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type TelegramAuthenticationStorage interface {
	Upsert(user initdata.User) (authentication.User, error)
}
