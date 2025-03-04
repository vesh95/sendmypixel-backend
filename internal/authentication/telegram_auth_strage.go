package authentication

import (
	"backend/pkg/authentication"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"strconv"
)

var InvalidCacheData = errors.New("invalid cache data")
var NotFoundInCache = errors.New("not found in cache")

const selectUserQuery = "SELECT u.id, u.username FROM telegram_users u WHERE u.id = $1"
const insertProviderData = "INSERT INTO telegram_users (id, username, lang) VALUES ($1, $2, $3) RETURNING id, username"

type TelegramAuthStorage struct {
	db    *sql.DB
	redis *redis.Client
	ctx   context.Context
}

func NewTelegramAuthStorage(db *sql.DB, redis *redis.Client) *TelegramAuthStorage {
	return &TelegramAuthStorage{db, redis, context.Background()}
}

const userCacheKeyPrefix = "tg_user_cache:"

func (s *TelegramAuthStorage) Upsert(user initdata.User) (authentication.User, error) {
	userData, err := s.tryCache(user.ID)
	if err == nil {
		return userData, nil
	} else if errors.Is(err, NotFoundInCache) {
		userData, err = s.selectOrInsert(user)
	}

	return userData, err
}

func (s *TelegramAuthStorage) tryCache(id int64) (authentication.User, error) {
	jsonUser, err := s.redis.Get(s.ctx, userCacheKeyPrefix+strconv.FormatInt(id, 10)).Result()
	var authUser authentication.User
	if err == nil {
		err = json.Unmarshal([]byte(jsonUser), &authUser)
		if err != nil {
			return authUser, InvalidCacheData
		}
	}

	if !errors.Is(err, redis.Nil) {
		return authUser, NotFoundInCache
	} else {
		return authUser, err
	}
}

func (s *TelegramAuthStorage) selectOrInsert(user initdata.User) (authentication.User, error) {
	var authUser authentication.User
	row := s.db.QueryRowContext(s.ctx, selectUserQuery, user.ID)
	err := row.Scan(&authUser.Id, &authUser.Username)
	if !errors.Is(err, sql.ErrNoRows) {
		niRow := s.db.QueryRowContext(s.ctx, insertProviderData, user.ID, user.Username, user.LanguageCode)
		err = niRow.Scan(&authUser.Id, &authUser.Username)
		if err == nil {
			return authUser, nil
		}
	}
	return authUser, err
}
