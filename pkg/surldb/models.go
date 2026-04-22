package surldb

import (
	"time"

	"github.com/uptrace/bun"
)

type Url struct {
	bun.BaseModel `bun:"table:urls,alias:u"`

	ID          int64     `bun:",pk,autoincrement"`
	ShortCode   string    `bun:",notnull"`
	ExpiredAt   time.Time `bun:",notnull"`
	OriginalUrl string    `bun:",notnull"`
	CreatedBy   string    `bun:",default:'annonymous'"`
}
