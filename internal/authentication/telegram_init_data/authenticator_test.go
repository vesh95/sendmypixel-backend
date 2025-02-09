package telegram_init_data

import (
	"backend/http/authentication"
	"errors"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"math/rand"
	"testing"
)

// testSecretKey тестовый ключ из документации
const testSecretKey = "7342037359:AAHI25ES9xCOMPokpYoz-p8XVrZUdygo2J4"

type mockStorage struct {
	initDataUser initdata.User
	internalUser authentication.User
}

type mockedTelegramAuthenticationStorage struct {
	users map[int64]mockStorage
}

func makeMockedTelegramAuthenticationStorage() *mockedTelegramAuthenticationStorage {
	users := make(map[int64]mockStorage)

	return &mockedTelegramAuthenticationStorage{users}
}

func (s *mockedTelegramAuthenticationStorage) getProviderUser(id int64) (initdata.User, bool) {
	for _, item := range s.users {
		if item.internalUser.Id == id {
			return item.initDataUser, true
		}
	}
	return initdata.User{}, false
}

func (s *mockedTelegramAuthenticationStorage) Upsert(initDataUser initdata.User) (authentication.User, error) {
	dbUser, ok := s.users[initDataUser.ID]
	if !ok {
		dbUser.initDataUser = initDataUser
		dbUser.internalUser = authentication.User{
			Id:       int64(rand.Int()),
			Username: initDataUser.Username,
		}
	}

	s.users[initDataUser.ID] = dbUser

	return dbUser.internalUser, nil
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name       string
		initData   string
		resultErr  error
		resultData initdata.User
	}{
		{
			name:      "TestValidDataFromDocumentation",
			initData:  "user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%20%2B%20-%20%3F%20%5C%2F%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F4FPEE4tmP3ATHa57u6MqTDih13LTOiMoKoLDRG4PnSA.svg%22%7D&chat_instance=8134722200314281151&chat_type=private&auth_date=1733509682&signature=TYJxVcisqbWjtodPepiJ6ghziUL94-KNpG8Pau-X7oNNLNBM72APCpi_RKiUlBvcqo5L-LAxIc3dnTzcZX_PDg&hash=a433d8f9847bd6addcc563bff7cc82c89e97ea0d90c11fe5729cae6796a36d73",
			resultErr: nil,
			resultData: initdata.User{
				AddedToAttachmentMenu: true,
				AllowsWriteToPm:       true,
				FirstName:             "Vladislav + - ? /",
				ID:                    279058397,
				IsBot:                 false,
				IsPremium:             true,
				LastName:              "Kibenko",
				Username:              "vdkfrost",
				LanguageCode:          "ru",
				PhotoURL:              "https://t.me/i/userpic/320/4FPEE4tmP3ATHa57u6MqTDih13LTOiMoKoLDRG4PnSA.svg",
			},
		},
		{
			name:       "TestInvalidSignature",
			initData:   "user=%7B%22id%22%3A379058397%2C%22first_name%22%3A%22Vladislav%20%2B%20-%20%3F%20%5C%2F%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F4FPEE4tmP3ATHa57u6MqTDih13LTOiMoKoLDRG4PnSA.svg%22%7D&chat_instance=8134722200314281151&chat_type=private&auth_date=1733509682&signature=TYJxVcisqbWjtodPepiJ6ghziUL94-KNpG8Pau-X7oNNLNBM72APCpi_RKiUlBvcqo5L-LAxIc3dnTzcZX_PDg&hash=a433d8f9847bd6addcc563bff7cc82c89e97ea0d90c11fe5729cae6796a36d73",
			resultErr:  authentication.InvalidSignatureError,
			resultData: initdata.User{},
		},
		{
			name:       "TestValidData",
			initData:   "user%3D%22%3A279058397%2C%22first_name%22%3A%22Vladislav%20%2B%20-%20%3F%20%5C%2F%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F4FPEE4tmP3ATHa57u6MqTDih13LTOiMoKoLDRG4PnSA.svg%22%7D%26chat_UlBvcqo5L-LAxIc3dnTzcZX_PDg%26hash%3Da433d8f9847bd6addcc563bff7cc82c89e97ea0d90c11fe5729cae6796a36d73",
			resultErr:  authentication.InvalidSignatureError,
			resultData: initdata.User{},
		},
	}

	storage := makeMockedTelegramAuthenticationStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authenticator := NewTelegramMiniAppAuthenticator(testSecretKey, storage)
			user, err := authenticator.Authenticate(test.initData)
			if !errors.Is(err, test.resultErr) {
				t.Errorf("Unexpected error: %s, except: %s", err, test.resultErr)
			}

			providerUser, _ := storage.getProviderUser(user.Id)
			if !compareResult(providerUser, test.resultData) {
				t.Errorf("Recieved: %+v\nexcept: %+v", providerUser, test.resultData)
			}
		})
	}
}

func compareResult(result, await initdata.User) bool {
	return result.ID == await.ID &&
		result.FirstName == await.FirstName &&
		result.LastName == await.LastName &&
		result.Username == await.Username
}
