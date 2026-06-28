// Package listing_client — тонкий HTTP-клиент к listing-service.
//
// Единственное, для чего messenger-service ходит в другой сервис синхронно:
// узнать владельца (продавца) объявления при создании треда переписки.
// Это происходит один раз на создание треда, не на каждое сообщение —
// дальше продавец/покупатель уже сохранены в самой messenger-service.
package listing_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	core_errors "messenger-service/internal/core/errors"
	core_logger "messenger-service/internal/core/logger"
)

type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	// Budget — общий потолок времени на все попытки и паузы между ними
	// вместе. 0 означает "без общего бюджета" (ограничены только
	// MaxRetries/MaxBackoff по отдельности).
	Budget time.Duration
}

type Client struct {
	baseURL string
	http    *http.Client
	retry   RetryConfig
	log     *core_logger.Logger
}

func NewClient(baseURL string, retry RetryConfig, log *core_logger.Logger) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 5 * time.Second},
		retry:   retry,
		log:     log,
	}
}

// Listing — то немногое, что нужно от ответа listing-service.
type Listing struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

// GetListing вызывает публичный GET /api/listings/:id (без авторизации),
// с экспоненциальным backoff-ретраем (+jitter) на сетевые/5xx-сбои.
// MaxRetries=0 означает "одна попытка, без повторов", а не "не пытаться
// вовсе" — нижняя граница цикла поэтому inclusive ("<="), не "<".
func (c *Client) GetListing(ctx context.Context, id uuid.UUID) (Listing, error) {
	if c.retry.Budget > 0 {
		budgetCtx, cancel := context.WithTimeout(ctx, c.retry.Budget)
		defer cancel()
		ctx = budgetCtx
	}

	var lastErr error

	for attempt := 0; attempt <= c.retry.MaxRetries; attempt++ {
		listing, err := c.getListing(ctx, id)
		if err == nil {
			return listing, nil
		}

		if errors.Is(err, core_errors.ErrNotFound) {
			return Listing{}, err
		}

		lastErr = err

		if attempt < c.retry.MaxRetries {
			backoff := c.retry.InitialBackoff * time.Duration(math.Pow(2, float64(attempt)))
			backoff = min(backoff, c.retry.MaxBackoff)
			backoff = fullJitter(backoff)

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return Listing{}, fmt.Errorf("retry budget exceeded: %w", ctx.Err())
			}
		}
	}

	if ctx.Err() != nil {
		return Listing{}, fmt.Errorf("retry budget exceeded: %w", ctx.Err())
	}

	return Listing{}, fmt.Errorf("failed after %d attempts: %w", c.retry.MaxRetries+1, lastErr)
}

// fullJitter возвращает случайную длительность в [0, d) — рассеивает
// повторы во времени, чтобы много клиентов, упёршихся в сбой одновременно,
// не били по восстанавливающемуся сервису синхронными залпами.
func fullJitter(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	return time.Duration(rand.Int64N(int64(d)))
}

func (c *Client) getListing(ctx context.Context, id uuid.UUID) (Listing, error) {
	url := fmt.Sprintf("%s/api/listings/%s", c.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.log.Error("build request:", zap.Error(err))
		return Listing{}, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.log.Error("call listing-service", zap.Error(err))
		return Listing{}, fmt.Errorf("call listing-service: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		c.log.Info("listing not found")
		return Listing{}, core_errors.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		c.log.Info("listing-service returned an unexpected status code")
		return Listing{}, fmt.Errorf("listing-service returned status %d", resp.StatusCode)
	}

	var listing Listing
	if err := json.NewDecoder(resp.Body).Decode(&listing); err != nil {
		return Listing{}, fmt.Errorf("decode response: %w", err)
	}
	return listing, nil
}
