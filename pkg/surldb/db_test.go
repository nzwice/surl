package surldb

import (
	"context"
	"fmt"
	"testing"

	"github.com/nzwice/surl/pkg/config"
)

func TestGetUrls(t *testing.T) {
	db, err := New(config.DBConfig{
		DSN: "postgres://surl:surl@localhost:5432/surl?sslmode=disable",
	}, true)
	if err != nil {
		panic(err)
	}
	var ctx = context.Background()
	var selectedUrl Url
	err = db.NewSelect().Model(&selectedUrl).Where("short_code = ?", "abc1234").Scan(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(selectedUrl)
}
