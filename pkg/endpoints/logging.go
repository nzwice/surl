package endpoints

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-kit/kit/endpoint"
)

func loggingMw() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				slog.InfoContext(
					ctx,
					"endpoint log",
					slog.Any("duration", time.Since(begin).Seconds()),
				)
			}(time.Now())
			return next(ctx, request)
		}
	}
}
