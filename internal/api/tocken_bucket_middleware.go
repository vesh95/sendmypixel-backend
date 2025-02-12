package api

import (
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity int
	tokens   int
	mutex    sync.Mutex
}

// NewTokenBucket - конструктор мидлвара, который создаёт структуру TokenBucket.
// capacity определяет кол-во токенов в бакете, а initialTokens сколько их будет изначально.
func NewTokenBucket(capacity, initialTokens int) *TokenBucket {
	tb := &TokenBucket{
		capacity: capacity,
		tokens:   initialTokens,
	}
	tb.startTokenRefill()

	return tb
}

// startTokenRefill запускает периодическое восстановление токенов в бакете.
func (tb *TokenBucket) startTokenRefill() {
	go func() {
		for {
			time.Sleep(time.Second)
			tb.mutex.Lock()
			if tb.tokens < tb.capacity {
				tb.tokens++
			}
			tb.mutex.Unlock()
		}
	}()
}

// consumeToken задействует токен из бакета.
func (tb *TokenBucket) consumeToken() {
	tb.mutex.Lock()
	for tb.tokens == 0 {
		tb.mutex.Unlock()
		time.Sleep(time.Millisecond * 100) // Настройка времени ожидания
		tb.mutex.Lock()
	}
	tb.tokens--
	tb.mutex.Unlock()
}

// RateLimitMiddleware это промежуточная функция обработчика http запросов для получения токена при выполнении запроса.
func (tb *TokenBucket) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tb.consumeToken()
		if tb.tokens == 0 {
			http.Error(w, "Token bucket limit reached", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
