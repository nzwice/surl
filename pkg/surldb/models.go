package surldb

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type Url struct {
	bun.BaseModel `bun:"table:urls,alias:u"`

	ID          int64     `bun:",pk,autoincrement"`
	ShortCode   string    `bun:",notnull"`
	ExpiredAt   time.Time `bun:",notnull"`
	OriginalUrl string    `bun:",notnull"`
	CreatedBy   string    `bun:",default:'annonymous'"`
	CreatedAt   time.Time `bun:",notnull"`
	UpdatedAt   time.Time `bun:",notnull"`
}

// BeforeAppendModel implements [schema.BeforeAppendModelHook].
func (u *Url) BeforeAppendModel(ctx context.Context, query schema.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		u.CreatedAt = time.Now().UTC()
		u.UpdatedAt = time.Now().UTC()
	case *bun.UpdateQuery:
		u.UpdatedAt = time.Now().UTC()
	}
	return nil
}

var _ bun.BeforeAppendModelHook = new(Url)
