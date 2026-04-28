package validateutils

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

func Struct(ctx context.Context, v any) error {
	if typ := reflect.TypeOf(v); typ == nil || typ.Kind() != reflect.Struct {
		slog.ErrorContext(ctx, "using struct validation on non-struct types", slog.Any("type", typ))
	}
	return validate.Struct(v)
}
