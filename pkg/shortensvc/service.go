package shortensvc

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"time"

	"github.com/nzwice/surl/pkg/surldb"
	"github.com/uptrace/bun"
)

var (
	ErrShortCodeNotFound  = errors.New("short code not found")
	ErrInvalidOriginalUrl = errors.New("invalid original url")
	ErrDBError            = errors.New("db error")
)

type Service interface {
	ShortenUrl(ctx context.Context, originalUrl string, alias *string, expiredAt *time.Time) (string, error)
	GetOriginalUrl(ctx context.Context, shortCode string) (string, error)
}

type Proxy func(Service) Service

type service struct {
	db *bun.DB
}

func New(db *bun.DB) Service {
	return &service{
		db: db,
	}
}

// GetOriginalUrl implements ShortenUrlService.
func (s *service) GetOriginalUrl(ctx context.Context, shortCode string) (string, error) {
	var existingUrl surldb.Url
	err := s.db.NewSelect().Model(&existingUrl).Where("short_code = ?", shortCode).Where("expired_at > now()").Column("original_url").Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrShortCodeNotFound
		}
		return "", errors.Join(ErrDBError, err)
	}
	return existingUrl.OriginalUrl, nil
}

// ShortenUrl implements ShortenUrlService.
func (s *service) ShortenUrl(ctx context.Context, originalUrl string, alias *string, expiredAt *time.Time) (string, error) {
	parsedUrl, err := url.Parse(originalUrl)
	if err != nil {
		return "", errors.Join(ErrInvalidOriginalUrl, err)
	}
	originalUrl = parsedUrl.String()
	var isNewAlias bool
	if alias != nil {
		isExisting, err := s.db.NewSelect().Table("urls").Where("short_code = ?", *alias).Exists(ctx)
		if err == nil {
			isNewAlias = !isExisting
		} else {
			return "", errors.Join(ErrDBError, err)
		}
	}
	var shortCode string
	if isNewAlias {
		shortCode = *alias
	} else {
		shortCode = s.generateShortCode(originalUrl)
	}
	var newUrl = surldb.Url{
		OriginalUrl: originalUrl,
		ShortCode:   shortCode,
		ExpiredAt:   s.expireOrDefault(expiredAt),
	}
	_, err = s.db.NewInsert().Model(&newUrl).Exec(ctx)
	if err != nil {
		return "", errors.Join(ErrDBError, err)
	}
	return newUrl.ShortCode, nil
}

var (
	shortCodeCharset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func (s *service) generateShortCode(_ string) string {
	nowNano := int(time.Now().UTC().UnixNano())
	base62 := func(n int) string {
		var sz = len(shortCodeCharset)
		var b []byte
		for n > 0 {
			b = append(b, shortCodeCharset[n%sz])
			n /= sz
		}
		for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
			b[i], b[j] = b[j], b[i]
		}
		return string(b)
	}
	return base62(nowNano)
}

func (s *service) expireOrDefault(expiredAt *time.Time) time.Time {
	if expiredAt != nil && time.Now().UTC().Before(*expiredAt) {
		return *expiredAt
	}
	return time.Now().UTC().Add(48 * time.Hour)
}
