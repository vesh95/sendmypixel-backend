package ratelimit

import (
	"backend/internal/api"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimitMiddlewareWithOnceRequest(t *testing.T) {
	tb := api.NewTokenBucket(1, 1, 1*time.Second)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "Hello, World!") })
	middleware := tb.RateLimitMiddleware(handler)

	// Создание тестового запроса
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRateLimitMiddlewareWith(t *testing.T) {
	// Инициализация ведра токенов с емкостью 1 и начальным количеством токенов 1
	tb := api.NewTokenBucket(1, 1, 1*time.Second)

	// Создание обработчика для тестирования
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Создание middleware с использованием ведра токенов
	middleware := tb.RateLimitMiddleware(handler)

	// Создание тестового запроса
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	rr = httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRateLimitMiddlewareWithRefill(t *testing.T) {
	// Инициализация ведра токенов с емкостью 1 и начальным количеством токенов 1
	tb := api.NewTokenBucket(1, 1, 1*time.Second)
	tb.StartTokenRefill()

	// Создание обработчика для тестирования
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Создание middleware с использованием ведра токенов
	middleware := tb.RateLimitMiddleware(handler)

	// Создание тестового запроса
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	rr = httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusTooManyRequests {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTooManyRequests)
	}

	time.Sleep(1 * time.Second)
	rr = httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	rr = httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusTooManyRequests {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTooManyRequests)
	}
}
