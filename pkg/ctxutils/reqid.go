package ctxutils

import "context"

const (
	reqIdKey = "requestId"
)

func GetRequestId(ctx context.Context) string {
	if v, ok := ctx.Value(reqIdKey).(string); ok {
		return v
	}
	return ""
}

func WithRequestId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, reqIdKey, id)
}
