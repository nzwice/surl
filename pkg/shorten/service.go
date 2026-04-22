package shorten

import (
	"context"
	"database/sql"
	"errors"
	"hash/fnv"
	"net/url"
	"strconv"
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

type service struct {
	db    *bun.DB
	inmem map[string]string
}

func New(db *bun.DB) Service {
	return &service{
		db:    db,
		inmem: map[string]string{},
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
		isExisting, err := s.db.NewSelect().Where("short_code = ?", *alias).Exists(ctx)
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

func (s *service) generateShortCode(originalUrl string) string {
	hash := fnv.New64a()
	hash.Write([]byte(originalUrl))
	calc := hash.Sum64()
	return strconv.FormatUint(calc, 10)
}

func (s *service) expireOrDefault(expiredAt *time.Time) time.Time {
	if expiredAt != nil && time.Now().UTC().Before(*expiredAt) {
		return *expiredAt
	}
	return time.Now().UTC().Add(48 * time.Hour)
}
