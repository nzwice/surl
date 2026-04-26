package shortensvc

import (
	"context"
	"fmt"
	"time"

	"github.com/nzwice/surl/pkg/kvstore"
)

var (
	urlCacheTTL = 5 * time.Minute
)

type cachedUrl struct {
	OriginalUrl string `json:"original_url"`
}

type cachedUrlSvc struct {
	client kvstore.Client
	next   Service
}

// GetOriginalUrl implements [Service].
func (c *cachedUrlSvc) GetOriginalUrl(ctx context.Context, shortCode string) (string, error) {
	cachedKey := c.buildUrlCacheKey(shortCode)
	var cachedObj cachedUrl
	err := c.client.DoGet(ctx, cachedKey, &cachedObj)
	if err == nil {
		return cachedObj.OriginalUrl, nil
	}
	originalUrl, err := c.next.GetOriginalUrl(ctx, shortCode)
	if err != nil {
		return originalUrl, err
	}
	_ = c.client.SetTTL(
		ctx,
		cachedKey,
		cachedUrl{
			OriginalUrl: originalUrl,
		},
		urlCacheTTL,
	)
	return originalUrl, nil
}

func (c *cachedUrlSvc) buildUrlCacheKey(shortCode string) string {
	return fmt.Sprintf("short_code:%s", shortCode)
}

// ShortenUrl implements [Service].
func (c *cachedUrlSvc) ShortenUrl(ctx context.Context, originalUrl string, alias *string, expiredAt *time.Time) (string, error) {
	shortCode, err := c.next.ShortenUrl(ctx, originalUrl, alias, expiredAt)
	if err != nil {
		return shortCode, err
	}
	cachedKey := c.buildUrlCacheKey(shortCode)
	_ = c.client.SetTTL(
		ctx,
		cachedKey,
		cachedUrl{
			OriginalUrl: originalUrl,
		},
		urlCacheTTL,
	)
	return shortCode, nil
}

func NewCache(client kvstore.Client) Proxy {
	return func(next Service) Service {
		return &cachedUrlSvc{
			next:   next,
			client: client,
		}
	}
}
