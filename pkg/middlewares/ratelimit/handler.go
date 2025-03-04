package ratelimit

import (
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity       int
	tokens         int
	refillDuration time.Duration
	mutex          sync.Mutex
}

// NewTokenBucket - конструктор мидлвара, который создаёт структуру TokenBucket.
// capacity определяет кол-во токенов в бакете, а initialTokens сколько их будет изначально.
// duration определяет через какой промежуток времени обновляется токен
func NewTokenBucket(capacity, initialTokens int, duration time.Duration) *TokenBucket {
	tb := &TokenBucket{
		capacity:       capacity,
		tokens:         initialTokens,
		refillDuration: duration,
	}

	return tb
}

// StartTokenRefill запускает периодическое восстановление токенов в бакете.
func (tb *TokenBucket) StartTokenRefill() {
	go func() {
		for {
			time.Sleep(tb.refillDuration)
			tb.mutex.Lock()
			if tb.tokens < tb.capacity {
				tb.tokens++
			}
			tb.mutex.Unlock()
		}
	}()
}

// Take задействует токен из бакета.
func (tb *TokenBucket) Take() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	if tb.tokens == 0 {
		return false
	}
	tb.tokens--

	return true
}

// Chain это промежуточная функция обработчика http запросов для получения токена при выполнении запроса.
func (tb *TokenBucket) Chain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tb.Take() {
			http.Error(w, "Token bucket limit reached", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
