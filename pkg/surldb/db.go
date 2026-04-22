package surldb

import (
	"database/sql"
	"fmt"

	"github.com/nzwice/surl/pkg/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func New(cfg config.DBConfig, debug bool) (*bun.DB, error) {

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(cfg.DSN),
	))

	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqldb.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqldb.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := sqldb.Ping(); err != nil {
		return nil, fmt.Errorf("fail to ping db: %v", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	if debug {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
		))

	}

	return db, nil
}
